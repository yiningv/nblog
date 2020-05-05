package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
)

// 标签列表
func (d *Dao) GetTags(pn, ps int) (tag []*model.Tag, pager *model.Page, err error) {
	db := d.DB.Table(model.TagTable)
	pager = &model.Page{
		Pn: pn,
		Ps: ps,
	}
	if err = db.Count(&pager.Total).Error; err != nil {
		log.Error(fmt.Sprintf("GetTags Count Error %v", err))
		return
	}
	if err = db.Order("article_count DESC").Offset((pn - 1) * ps).Limit(ps).Find(&tag).Error; err != nil {
		log.Error(fmt.Sprintf("GetTags Error %v", err))
	}
	return
}

// 添加标签
func (d *Dao) AddTag(arg *model.Tag) error {
	return d.DB.Create(arg).Error
}

// 删除标签的同时需要删除标签的关联
func (d *Dao) DeleteTag(id int64) (err error) {
	err = d.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.TagTable).Where("id=?", id).Delete(&model.Tag{}).Error; err != nil {
			return
		}
		err = tx.Table(model.TagRefTable).Where("tag_id=?", id).Delete(&model.TagRef{}).Error
		return
	})
	return
}

// 批量删除标签
func (d *Dao) BatchDeleteTag(ids []int64) (err error) {
	err = d.DB.Transaction(func(tx *gorm.DB) (err error) {
		if err = tx.Table(model.TagTable).Where("id IN (?)", ids).Delete(&model.Tag{}).Error; err != nil {
			return
		}
		err = tx.Table(model.TagRefTable).Where("tag_id IN (?)", ids).Delete(&model.TagRef{}).Error
		return
	})
	return
}

// 修改标签
func (d *Dao) UpdateTag(arg *model.Tag) error {
	return d.DB.Table(model.TagTable).Update(arg).Error
}
