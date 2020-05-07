package service

import (
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/model"
)

// 获取站点配置
func (srv *Service) GetSiteConfig() (siteConf []*model.SiteConfig, err error) {
	err = srv.dao.Table(model.SiteConfigTable).Find(&siteConf).Error
	return
}

// 添加站点配置
func (srv *Service) AddSiteConfig(arg *model.SiteConfig) (err error) {
	return srv.dao.Table(model.SiteConfigTable).Create(arg).Error
}

// 删除站点配置
func (srv *Service) DeleteSiteConfig(id int64) (err error) {
	return srv.dao.Table(model.SiteConfigTable).Where("id=?", id).Delete(&model.SiteConfig{}).Error
}

// 批量删除站点配置
func (srv *Service) BatchDeleteSiteConfig(ids []int64) (err error) {
	return srv.dao.Table(model.SiteConfigTable).Where("id IN (?)", ids).Delete(&model.SiteConfig{}).Error
}

// 修改站点配置
func (srv *Service) UpdateSiteConfig(arg *model.SiteConfig) (err error) {
	return srv.dao.Table(model.SiteConfigTable).Where("id=?", arg.ID).Update(arg).Error
}

// 修改站点配置
func (srv *Service) BatchUpdateSiteConfig(fromDB []*model.SiteConfig, confMap map[string]*model.SiteConfig) (err error) {
	save := make([]*model.SiteConfig, 0)
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
			if err = tx.Table(model.SiteConfigTable).Where("id IN (?)", delIds).Delete(&model.SiteConfig{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			config := save[i]
			if err = tx.Table(model.SiteConfigTable).Save(config).Error; err != nil {
				return
			}
		}
		return
	})
	return
}
