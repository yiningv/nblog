package service

import (
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
)

// 获取资源配置
func (srv *Service) GetSourceConfig() (siteConf []*model.SourceConfig, err error) {
	err = srv.dao.Table(model.SourceConfigTable).Find(&siteConf).Error
	return
}

// 添加资源配置
func (srv *Service) AddSourceConfig(arg *model.SourceConfig) (err error) {
	return srv.dao.Table(model.SourceConfigTable).Create(arg).Error
}

// 删除资源配置
func (srv *Service) DeleteSourceConfig(id int64) (err error) {
	return srv.dao.Table(model.SourceConfigTable).Where("id=?", id).Delete(&model.SourceConfig{}).Error
}

// 批量删除资源配置
func (srv *Service) BatchDeleteSourceConfig(ids []int64) (err error) {
	return srv.dao.Table(model.SourceConfigTable).Where("id IN (?)", ids).Delete(&model.SourceConfig{}).Error
}

// 修改资源配置
func (srv *Service) UpdateSourceConfig(arg *model.SourceConfig) (err error) {
	return srv.dao.Table(model.SourceConfigTable).Where("id=?", arg.ID).Update(arg).Error
}

// 修改资源配置
func (srv *Service) BatchUpdateSourceConfig(fromDB []*model.SourceConfig, confMap map[string]*model.SourceConfig) (err error) {
	save := make([]*model.SourceConfig, 0)
	delIds := make([]int, 0)
	for i := range fromDB {
		scDB := fromDB[i]
		if sc, ok := confMap[scDB.Name]; ok {
			sc.ID = scDB.ID
			if sc.LastEditedTime < scDB.LastEditedTime {
				save = append(save, sc)
			}
			delete(confMap, scDB.Name)
		} else {
			delIds = append(delIds, scDB.ID)
		}
	}
	for k := range confMap {
		save = append(save, confMap[k])
	}
	if len(save) == 0 && len(delIds) == 0 {
		return
	}
	err = srv.dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Table(model.SourceConfigTable).Where("id IN (?)", delIds).Delete(&model.SourceConfig{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			config := save[i]
			if err = tx.Table(model.SourceConfigTable).Save(config).Error; err != nil {
				return
			}
		}
		return
	})
	return
}
