package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/cache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

var Archive = &archiveService{}

type archiveService struct{}

// 归档列表
func (srv *archiveService) GetArchive(pn, ps int) (pager *model.ArchivePager, err error) {
	pager = &model.ArchivePager{}
	page := &model.Page{
		Pn: pn,
		Ps: ps,
	}
	pager.Page = page
	if err = dao.Model(&model.Archive{}).
		Order("month DESC, month DESC").Count(&page.Total).
		Offset((pn - 1) * ps).Limit(ps).
		Find(&pager.Items).Error; err != nil {
		log.Error(fmt.Sprintf("GetArchive Error %v", err))
	}
	return
}

func (srv *archiveService) GetArchiveAll() (archives []*model.Archive, err error) {
	err = dao.Find(&archives).Error
	return
}

// 添加归档
func (srv *archiveService) AddTag(arg *model.Archive) error {
	return dao.Create(arg).Error
}

// 删除归档的同时需要删除标签的关联
func (srv *archiveService) DeleteArchive(id int) (err error) {
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("id=?", id).Delete(&model.Archive{}).Error; err != nil {
			return
		}
		err = tx.Where("id2=? AND type=?", id, model.CorrelationPostsArchive).Delete(&model.Correlation{}).Error
		return
	})
	return
}

// 批量删除归档
func (srv *archiveService) BatchDeleteArchive(ids []int) (err error) {
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("id IN (?)", ids).Delete(&model.Archive{}).Error; err != nil {
			return
		}
		err = tx.Where("id2 IN (?) AND type=?", ids, model.CorrelationPostsArchive).Delete(&model.Correlation{}).Error
		return
	})
	return
}

// 修改归档
func (srv *archiveService) UpdateTag(arg *model.Archive) error {
	return dao.Update(arg).Error
}

func (srv *archiveService) GetPostsArchive() (arg []*model.Correlation, err error) {
	err = dao.Where("type=?", model.CorrelationPostsArchive).Find(&arg).Error
	return
}

func (srv *archiveService) loadArchiveCache() (err error) {
	postsAll := cache.Posts.GetAll()
	var archiveAll []*model.Archive
	archiveAll, err = srv.GetArchiveAll()
	if err != nil {
		return
	}
	var postsArchive []*model.Correlation
	postsArchive, err = srv.GetPostsArchive()
	if err != nil {
		return
	}
	archiveMap := make(map[string]*model.Archive)
	for i := range archiveAll {
		archive := archiveAll[i]
		archiveMap[archive.Slug] = archive
	}
	newData := make(map[string]*model.ArchivePosts)
	for _, p_a := range postsArchive {
		pageId := p_a.Str1
		slug := p_a.Str2
		if postsAll[pageId] == nil {
			continue
		}
		if archivePosts, ok := newData[slug]; ok {
			archivePosts.SortPosts = append(archivePosts.SortPosts, postsAll[pageId])
		} else {
			newData[slug] = &model.ArchivePosts{
				Archive:   archiveMap[slug],
				SortPosts: []*model.Posts{postsAll[pageId]},
			}
		}
	}
	cache.Archive.Replace(newData)
	return
}

func (srv *archiveService) syncArchive() (err error) {
	all := cache.Posts.GetAll()
	archiveUpdate := make(map[string]*model.ArchivePosts)
	for i := range all {
		posts := all[i]
		pTime := posts.PTime
		slug := fmt.Sprintf("%d-%d", pTime.Year(), pTime.Month())
		archive, ok := archiveUpdate[slug]
		if ok {
			archive.SortPosts = append(archive.SortPosts, posts)
			archive.Archive.PostsCount++
		} else {
			archiveUpdate[slug] = &model.ArchivePosts{
				Archive: &model.Archive{
					PostsCount: 1,
				},
				SortPosts: []*model.Posts{posts},
			}
		}
		archiveUpdate[slug].Archive.Year = fmt.Sprintf("%d", pTime.Year())
		archiveUpdate[slug].Archive.Month = fmt.Sprintf("%d", pTime.Month())
		archiveUpdate[slug].Archive.Slug = slug
	}

	archiveAll := cache.Archive.GetAll()

	save := make([]*model.Archive, 0)
	delIds := make([]int, 0)
	for s := range archiveUpdate {
		aUpdate := archiveUpdate[s]
		if aCache, ok := archiveAll[s]; ok {
			aUpdate.Archive.ID = aCache.Archive.ID
			// tag对应的文章数量有变化, 需要更新
			if aCache.Archive.PostsCount != aUpdate.Archive.PostsCount {
				save = append(save, aUpdate.Archive)
			}
			delete(archiveAll, s)
		} else {
			save = append(save, aUpdate.Archive)
		}
	}
	// 缓存中剩下的是需要删除的数据
	for _, tCache := range archiveAll {
		delIds = append(delIds, tCache.Archive.ID)
	}

	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?)", delIds).Delete(&model.Archive{}).Error; err != nil {
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

	var postsArchive []*model.Correlation
	postsArchive, err = srv.GetPostsArchive()
	if err != nil {
		return
	}
	paMap := make(map[string]*model.Correlation)
	for i := range postsArchive {
		pt := postsArchive[i]
		paMap[pt.Str1+pt.Str2] = pt
	}

	savePA := make([]*model.Correlation, 0)
	delIdsPA := make([]int, 0)
	for s := range archiveUpdate {
		aUpdate := archiveUpdate[s]
		archive := aUpdate.Archive
		posts := aUpdate.SortPosts
		for _, p := range posts {
			paKey := p.PageId + archive.Slug
			if pt, ok := paMap[paKey]; ok {
				// posts和archive的id重新生成过
				if pt.ID1 != p.ID || pt.ID2 != archive.ID {
					pt.ID1 = p.ID
					pt.ID2 = archive.ID
					savePA = append(savePA, pt)
				}
				delete(paMap, paKey)
			} else {
				// 新的关联关系
				paNew := &model.Correlation{
					ID1:  p.ID,
					ID2:  archive.ID,
					Str1: p.PageId,
					Str2: archive.Slug,
					Type: model.CorrelationPostsArchive,
				}
				savePA = append(savePA, paNew)
			}
		}
	}

	for _, pt := range paMap {
		delIdsPA = append(delIdsPA, pt.ID)
	}

	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?) AND type=?", delIdsPA, model.CorrelationPostsArchive).Delete(&model.Correlation{}).Error; err != nil {
				return
			}
		}

		for i := range savePA {
			pt := savePA[i]
			if err = tx.Save(pt).Error; err != nil {
				return
			}
		}
		return
	})

	cache.Archive.Replace(archiveUpdate)

	return
}
