package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

// 分类列表
func (srv *Service) GetCategories() (cates []*model.Category, err error) {
	if err = srv.dao.Table(model.CategoryTable).Order("order_num DESC").Find(&cates).Error; err != nil {
		log.Error(fmt.Sprintf("GetCategories Error %v", err))
	}
	return
}

// 添加分类
func (srv *Service) AddCategory(arg *model.Category) error {
	return srv.dao.Create(arg).Error
}

// 修改分类
func (srv *Service) UpdateCategory(arg *model.Category) error {
	return srv.dao.Table(model.CategoryTable).Update(arg).Error
}

// 删除分类
func (srv *Service) DeleteCategories(ids []int64) (err error) {
	switch len(ids) {
	case 0:
		return
	case 1:
		if err = srv.deleteCategory(ids[0]); err != nil {
			return
		}
	default:
		if err = srv.batchDeleteCategory(ids); err != nil {
			return
		}
	}
	return
}

// 删除分类的同时需要删除分类的关联
func (srv *Service) deleteCategory(id int64) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.CategoryTable).Where("id=?", id).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.CategoryRefTable).Where("category_id=?", id).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 批量删除分类
func (srv *Service) batchDeleteCategory(ids []int64) (err error) {
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.CategoryTable).Where("id IN (?)", ids).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.CategoryRefTable).Where("category_id IN (?)", ids).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}
