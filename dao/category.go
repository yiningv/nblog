package dao

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"github.com/yiningv/nblog/pub/util"
)

// 分类列表
func (d *Dao) GetCategories() (cates []*model.Category, err error) {
	if err = d.DB.Table(model.CategoryTable).Order("order_num DESC").Find(&cates).Error; err != nil {
		log.Error(fmt.Sprintf("GetCategories Error %v", err))
	}
	return
}

// 添加分类
func (d *Dao) AddCategory(arg *model.Category) error {
	return d.DB.Create(arg).Error
}

// 删除分类的同时需要删除分类的关联
func (d *Dao) DeleteCategory(id int64) (err error) {
	err = d.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.CategoryTable).Where("id=?", id).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.CategoryRefTable).Where("category_id=?", id).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 批量删除分类
func (d *Dao) BatchDeleteCategory(ids []int64) (err error) {
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
		if err = tx.Table(model.CategoryTable).Where(buf.String()).Delete(&model.Category{}).Error; err != nil {
			return
		}
		err = tx.Table(model.CategoryRefTable).Where(buf2.String()).Delete(&model.CategoryRef{}).Error
		return
	})
	return
}

// 修改分类
func (d *Dao) UpdateCategory(arg *model.Category) error {
	return d.DB.Table(model.CategoryTable).Update(arg).Error
}
