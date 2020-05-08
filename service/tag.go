package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

// 标签列表
//
func (srv *Service) GetTags(pn, ps int) (pager *model.TagPager, err error) {
	pager = &model.TagPager{}
	dao := srv.dao.Table(model.TagTable)
	page := &model.Page{
		Pn: pn,
		Ps: ps,
	}
	pager.Page = page
	if err = dao.Count(&page.Total).Error; err != nil {
		log.Error(fmt.Sprintf("GetTags Count Error %v", err))
		return
	}
	if err = dao.Order("article_count DESC").Offset((pn - 1) * ps).Limit(ps).Find(&pager.Items).Error; err != nil {
		log.Error(fmt.Sprintf("GetTags Error %v", err))
	}
	return
}

// 添加标签
func (srv *Service) AddTag(arg *model.Tag) error {
	return srv.dao.Create(arg).Error
}

// 删除标签的同时需要删除标签的关联
func (srv *Service) DeleteTag(id int64) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.TagTable).Where("id=?", id).Delete(&model.Tag{}).Error; err != nil {
			return
		}
		err = tx.Table(model.TagRefTable).Where("tag_id=?", id).Delete(&model.TagRef{}).Error
		return
	})
	return
}

// 批量删除标签
func (srv *Service) BatchDeleteTag(ids []int64) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.TagTable).Where("id IN (?)", ids).Delete(&model.Tag{}).Error; err != nil {
			return
		}
		err = tx.Table(model.TagRefTable).Where("tag_id IN (?)", ids).Delete(&model.TagRef{}).Error
		return
	})
	return
}

// 修改标签
func (srv *Service) UpdateTag(arg *model.Tag) error {
	return srv.dao.Table(model.TagTable).Update(arg).Error
}
