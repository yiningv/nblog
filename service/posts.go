package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

// 文章列表
func (srv *Service) GetPostsPager(pn, ps int) (pager *model.PostsPager, err error) {
	pager = &model.PostsPager{}
	dao := srv.dao.Table(model.PostsTable)
	page := &model.Page{
		Pn: pn,
		Ps: ps,
	}
	pager.Page = page
	if err = dao.Count(&page.Total).Error; err != nil {
		log.Error(fmt.Sprintf("GetPostsPager Count Error %v", err))
		return
	}
	var posts []*model.Posts
	pager.Items = posts
	if err = dao.Order("ctime DESC").Offset((pn - 1) * ps).Limit(ps).Find(&posts).Error; err != nil {
		log.Error(fmt.Sprintf("GetPostsPager Error %v", err))
	}
	return
}

// 根据ID获取文章信息
func (srv *Service) GetPosts(id int) (posts *model.Posts, err error) {
	err = srv.dao.Table(model.PostsTable).Where("id=?", id).First(posts).Error
	return
}

// 添加文章
func (srv *Service) AddPosts(arg *model.Posts) error {
	return srv.dao.Create(arg).Error
}

// 删除文章
func (srv *Service) DeletePosts(id int) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.PostsTable).Where("id=?", id).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.ArticleTable).Where("category_id=?", id).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 批量删除文章
func (srv *Service) BatchDeletePosts(ids []int64) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.ArticleTable).Where("id IN (?)", ids).Delete(&model.Article{}).Error; err != nil {
			return
		}
		err = tx.Table(model.CategoryRefTable).Where("category_id IN (?)", ids).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 修改文章
func (srv *Service) UpdatePosts(arg *model.Article) error {
	return srv.dao.Table(model.ArticleTable).Update(arg).Error
}
