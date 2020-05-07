package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

// 文章列表
func (srv *Service) GetArticles(pn, ps int) (pager *model.ArticlePager, err error) {
	pager = &model.ArticlePager{}
	dao := srv.dao.Table(model.ArticleTable)
	page := &model.Page{
		Pn: pn,
		Ps: ps,
	}
	pager.Page = page
	if err = dao.Count(&page.Total).Error; err != nil {
		log.Error(fmt.Sprintf("GetArticles Count Error %v", err))
		return
	}
	var arts []*model.Article
	pager.Items = arts
	if err = dao.Order("ctime DESC").Offset((pn - 1) * ps).Limit(ps).Find(&arts).Error; err != nil {
		log.Error(fmt.Sprintf("GetArticles Error %v", err))
	}
	return
}

// 根据ID获取文章信息
func (srv *Service) GetArticle(id int64) (art *model.Article, err error) {
	err = srv.dao.Table(model.ArticleTable).Where("id=?", id).First(art).Error
	return
}

// 添加文章
func (srv *Service) AddArticle(arg *model.Article) error {
	return srv.dao.Create(arg).Error
}

// 删除文章
func (srv *Service) DeleteArticle(id int64) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.ArticleTable).Where("id=?", id).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.ArticleTable).Where("category_id=?", id).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 批量删除文章
func (srv *Service) BatchDeleteArticle(ids []int64) (err error) {
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
func (srv *Service) UpdateArticle(arg *model.Article) error {
	return srv.dao.Table(model.ArticleTable).Update(arg).Error
}
