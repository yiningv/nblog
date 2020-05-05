package dao

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"github.com/yiningv/nblog/pub/util"
)

// 文章列表
func (d *Dao) GetArticles(pn, ps int) (arts []*model.Article, pager *model.Page, err error) {
	db := d.DB.Table(model.ArticleTable)
	pager = &model.Page{
		Pn: pn,
		Ps: ps,
	}
	if err = db.Count(&pager.Total).Error; err != nil {
		log.Error(fmt.Sprintf("GetArticles Count Error %v", err))
		return
	}
	if err = db.Order("ctime DESC").Offset((pn - 1) * ps).Limit(ps).Find(&arts).Error; err != nil {
		log.Error(fmt.Sprintf("GetArticles Error %v", err))
	}
	return
}

// 根据ID获取文章信息
func (d *Dao) GetArticle(id int64) (art *model.Article, err error) {
	err = d.DB.Table(model.ArticleTable).Where("id=?", id).First(art).Error
	return
}

// 添加文章
func (d *Dao) AddArticle(arg *model.Article) error {
	return d.DB.Create(arg).Error
}

// 删除文章
func (d *Dao) DeleteArticle(id int64) (err error) {
	err = d.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.ArticleTable).Where("id=?", id).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.ArticleTable).Where("category_id=?", id).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 批量删除文章
func (d *Dao) BatchDeleteArticle(ids []int64) (err error) {
	idStr := util.JoinInts(ids)
	var buf bytes.Buffer
	buf.WriteString("id IN (")
	buf.WriteString(idStr)
	buf.WriteString(")")

	var buf2 bytes.Buffer
	buf2.WriteString("category_id IN (")
	buf2.WriteString(idStr)
	buf2.WriteString(")")
	err = d.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.ArticleTable).Where(buf.String()).Delete(&model.Article{}).Error; err != nil {
			return
		}
		err = tx.Table(model.CategoryRefTable).Where(buf2.String()).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 修改文章
func (d *Dao) UpdateArticle(arg *model.Article) error {
	return d.DB.Table(model.ArticleTable).Update(arg).Error
}
