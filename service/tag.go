package service

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/cache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"strings"
)

var Tag = &tagService{}

type tagService struct{}

// 标签列表
func (srv *tagService) GetTags(pn, ps int) (pager *model.TagPager, err error) {
	pager = &model.TagPager{}
	page := &model.Page{
		Pn: pn,
		Ps: ps,
	}
	pager.Page = page
	if err = dao.Model(&model.Posts{}).
		Order("posts_count DESC").Count(&page.Total).
		Offset((pn - 1) * ps).Limit(ps).
		Find(&pager.Items).Error; err != nil {
		log.Error(fmt.Sprintf("GetTags Error %v", err))
	}
	return
}

func (srv *tagService) GetTagsAll() (tags []*model.Tag, err error) {
	err = dao.Find(&tags).Error
	return
}

// 添加标签
func (srv *tagService) AddTag(arg *model.Tag) error {
	return dao.Create(arg).Error
}

// 删除标签的同时需要删除标签的关联
func (srv *tagService) DeleteTag(id int) (err error) {
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("id=?", id).Delete(&model.Tag{}).Error; err != nil {
			return
		}
		err = tx.Where("id2=? AND type=?", id, model.CorrelationPostsTag).Delete(&model.Correlation{}).Error
		return
	})
	return
}

// 批量删除标签
func (srv *tagService) BatchDeleteTag(ids []int) (err error) {
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("id IN (?)", ids).Delete(&model.Tag{}).Error; err != nil {
			return
		}
		err = tx.Where("id2 IN (?) AND type=?", ids, model.CorrelationPostsTag).Delete(&model.Correlation{}).Error
		return
	})
	return
}

// 修改标签
func (srv *tagService) UpdateTag(arg *model.Tag) error {
	return dao.Update(arg).Error
}

func (srv *tagService) GetPostsTag() (arg []*model.Correlation, err error) {
	err = dao.Where("type=?", model.CorrelationPostsTag).Find(&arg).Error
	return
}

func (srv *tagService) loadTagCache() (err error) {
	postsAll := cache.Posts.GetAll()
	var tagAll []*model.Tag
	tagAll, err = srv.GetTagsAll()
	if err != nil {
		return
	}
	var postsTag []*model.Correlation
	postsTag, err = srv.GetPostsTag()
	if err != nil {
		return
	}
	tagMap := make(map[string]*model.Tag)
	for i := range tagAll {
		tag := tagAll[i]
		tagMap[tag.Slug] = tag
	}
	newData := make(map[string]*model.TagPosts)
	for _, p_t := range postsTag {
		pageId := p_t.Str1
		slug := p_t.Str2
		if postsAll[pageId] == nil {
			continue
		}
		if tagPosts, ok := newData[slug]; ok {
			tagPosts.SortPosts = append(tagPosts.SortPosts, postsAll[pageId])
		} else {
			newData[slug] = &model.TagPosts{
				Tag:       tagMap[slug],
				SortPosts: []*model.Posts{postsAll[pageId]},
			}
		}
	}
	cache.Tag.Replace(newData)
	return
}

func (srv *tagService) syncTag() (err error) {
	all := cache.Posts.GetAll()
	tagUpdate := make(map[string]*model.TagPosts)
	for i := range all {
		posts := all[i]
		tags := posts.Tags
		if tags == "" {
			continue
		}
		split := strings.Split(tags, ",")
		for _, tName := range split {
			tName = strings.TrimSpace(tName)
			s := slug.Make(tName)
			tag, ok := tagUpdate[s]
			if ok {
				tag.SortPosts = append(tag.SortPosts, posts)
				tag.Tag.PostsCount++
			} else {
				tagUpdate[s] = &model.TagPosts{
					Tag: &model.Tag{
						PostsCount: 1,
					},
					SortPosts: []*model.Posts{posts},
				}
			}
			tagUpdate[s].Tag.Name = tName
			tagUpdate[s].Tag.Slug = s
		}
	}

	tagCache := cache.Tag.GetAll()

	save := make([]*model.Tag, 0)
	delIds := make([]int, 0)
	for s := range tagUpdate {
		tUpdate := tagUpdate[s]
		if tCache, ok := tagCache[s]; ok {
			tUpdate.Tag.ID = tCache.Tag.ID
			// tag对应的文章数量有变化, 需要更新
			if tCache.Tag.PostsCount != tUpdate.Tag.PostsCount {
				save = append(save, tUpdate.Tag)
			}
			delete(tagCache, s)
		} else {
			save = append(save, tUpdate.Tag)
		}
	}
	// 缓存中剩下的是需要删除的数据
	for _, tCache := range tagCache {
		delIds = append(delIds, tCache.Tag.ID)
	}

	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?)", delIds).Delete(&model.Tag{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			tag := save[i]
			if err = tx.Save(tag).Error; err != nil {
				return
			}
		}
		return
	})

	var postsTag []*model.Correlation
	postsTag, err = srv.GetPostsTag()
	if err != nil {
		return
	}
	ptMap := make(map[string]*model.Correlation)
	for i := range postsTag {
		pt := postsTag[i]
		ptMap[pt.Str1+pt.Str2] = pt
	}

	savePT := make([]*model.Correlation, 0)
	delIdsPT := make([]int, 0)
	for s := range tagUpdate {
		tUpdate := tagUpdate[s]
		tag := tUpdate.Tag
		posts := tUpdate.SortPosts
		for _, p := range posts {
			ptKey := p.PageId + tag.Slug
			if pt, ok := ptMap[ptKey]; ok {
				// posts和tag的id重新生成过
				if pt.ID1 != p.ID || pt.ID2 != tag.ID {
					pt.ID1 = p.ID
					pt.ID2 = tag.ID
					savePT = append(savePT, pt)
				}
				delete(ptMap, ptKey)
			} else {
				// 新的关联关系
				ptNew := &model.Correlation{
					ID1:  p.ID,
					ID2:  tag.ID,
					Str1: p.PageId,
					Str2: tag.Slug,
					Type: model.CorrelationPostsTag,
				}
				savePT = append(savePT, ptNew)
			}
		}
	}

	for _, pt := range ptMap {
		delIdsPT = append(delIdsPT, pt.ID)
	}

	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?) AND type=?", delIdsPT, model.CorrelationPostsTag).Delete(&model.Correlation{}).Error; err != nil {
				return
			}
		}

		for i := range savePT {
			pt := savePT[i]
			if err = tx.Save(pt).Error; err != nil {
				return
			}
		}
		return
	})

	cache.Tag.Replace(tagUpdate)

	return
}
