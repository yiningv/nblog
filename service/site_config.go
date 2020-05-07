package service

import (
	"github.com/yiningv/nblog/model"
)

// 获取站点配置
func (srv *Service) GetSiteConfig() (siteConf []*model.SiteConfig, err error) {
	err = srv.dao.Table(model.SiteConfigTable).Find(&siteConf).Error
	return
}

// 添加站点配置
func (srv *Service) AddSiteConfig(arg *model.SiteConfig) error {
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
func (srv *Service) UpdateSiteConfig(arg *model.SiteConfig) error {
	args := make(map[string]interface{})
	args["name"] = arg.Name
	args["type"] = arg.Type
	args["value"] = arg.Value
	args["desc"] = arg.Desc
	return srv.dao.Table(model.SiteConfigTable).Where("id=?", arg.ID).Update(args).Error
}
