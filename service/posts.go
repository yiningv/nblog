package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/cache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/notion"
	"github.com/yiningv/nblog/pub/log"
	"os/exec"
	"path/filepath"
)

var Posts = &postsService{}

type postsService struct{}

// 文章列表
func (srv *postsService) GetPostsPager(pn, ps int) (pager *model.PostsPager, err error) {
	pager = &model.PostsPager{}
	page := &model.Page{
		Pn: pn,
		Ps: ps,
	}
	pager.Page = page
	var posts []*model.Posts
	pager.Items = posts
	if err = dao.Model(&model.Posts{}).
		Order("p_time DESC").Count(&page.Total).
		Offset((pn - 1) * ps).Limit(ps).
		Find(&posts).Error; err != nil {
		log.Error(fmt.Sprintf("GetPostsPager Error %v", err))
	}
	return
}

func (srv *postsService) GetPostsAll() (posts []*model.Posts, err error) {
	err = dao.Find(&posts).Error
	return
}

// 根据ID获取文章信息
func (srv *postsService) GetPosts(id int) (posts *model.Posts, err error) {
	err = dao.Where("id=?", id).First(posts).Error
	return
}

// 添加文章
func (srv *postsService) AddPosts(arg *model.Posts) error {
	return dao.Create(arg).Error
}

// 删除文章
func (srv *postsService) DeletePosts(id int) (err error) {
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("id=?", id).Delete(&model.Posts{}).Error; err != nil {
			return
		}
		err = tx.Where("id1=?", id).Delete(&model.Correlation{}).Error
		return
	})
	return
}

// 批量删除文章
func (srv *postsService) BatchDeletePosts(ids []int) (err error) {
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Where("id IN (?)", ids).Delete(&model.Posts{}).Error; err != nil {
			return
		}
		err = tx.Where("id1 IN (?)", ids).Delete(&model.Correlation{}).Error
		return
	})
	return
}

// 修改文章
func (srv *postsService) UpdatePosts(arg *model.Posts) error {
	return dao.Update(arg).Error
}

// 同步文章缓存
func (srv *postsService) loadPostsCache() (err error) {
	var postsAll []*model.Posts
	postsAll, err = srv.GetPostsAll()
	cache.Posts.Update(postsAll)
	return
}

// 同步文章信息
func (srv *postsService) syncPosts() (err error) {
	var postsUpdate map[string]*model.Posts
	postsUpdate, err = notion.GetPosts()
	if err != nil {
		log.Error(fmt.Sprintf("notion.GetPosts failed: %v", err))
		return
	}
	postsCache := cache.Posts.GetAll()

	newData := make([]*model.Posts, 0)
	save := make([]*model.Posts, 0)
	delIds := make([]int, 0)
	for name := range postsUpdate {
		pUpdate := postsUpdate[name]
		if pCache, ok := postsCache[name]; ok {
			pUpdate.ID = pCache.ID
			// 最后更新时间有变化或者顺序有变化时，需要对数据做更新
			if pUpdate.LastEditedTime != pCache.LastEditedTime {
				save = append(save, pUpdate)
			}
			delete(postsCache, name)
		} else {
			save = append(save, pUpdate)
		}
		newData = append(newData, pUpdate)
	}
	// 缓存中剩下的数据需要删除
	for _, sCache := range postsCache {
		delIds = append(delIds, sCache.ID)
	}

	if len(save) == 0 && len(delIds) == 0 {
		return
	}
	for i := range save {
		p := save[i]
		html, err := trace(p.PageId)
		if err != nil {
			p.Content = err.Error()
		} else {
			p.Content = html
		}
		p.Slug = slug.Make(p.Title)
	}

	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?)", delIds).Delete(&model.Posts{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			posts := save[i]
			if err = tx.Save(posts).Error; err != nil {
				return
			}
		}
		return
	})
	cache.Posts.Update(newData)
	return
}

func trace(pageId string) (html string, err error) {
	url := fmt.Sprintf("https://notion.so/%s", pageId)
	scriptPath := filepath.Join("puppeteer", "trace.js")
	cmd := exec.Command("node", scriptPath, url)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return
	}
	if stderr.Len() != 0 {
		err = errors.New("sync HTML error")
		return
	}
	html = stdout.String()
	return
}
