package dao

import (
	"github.com/yiningv/nblog/model"
)

// 获取站点配置
func (d *Dao) GetSiteConfig() (siteConf []*model.SiteConfig, err error) {
	err = d.DB.Table(model.SiteConfigTable).Find(&siteConf).Error
	return
}

// 添加站点配置
func (d *Dao) AddSiteConfig(arg *model.SiteConfig) error {
	return d.DB.Table(model.SiteConfigTable).Create(arg).Error
}

// 删除站点配置
func (d *Dao) DeleteSiteConfig(id int64) (err error) {
	return d.DB.Table(model.SiteConfigTable).Where("id=?", id).Delete(&model.SiteConfig{}).Error
}

// 批量删除站点配置
func (d *Dao) BatchDeleteSiteConfig(ids []int64) (err error) {
	return d.DB.Table(model.SiteConfigTable).Where("id IN (?)", ids).Delete(&model.SiteConfig{}).Error
}

// 修改站点配置
func (d *Dao) UpdateSiteConfig(arg *model.SiteConfig) error {
	args := make(map[string]interface{})
	args["name"] = arg.Name
	args["type"] = arg.Type
	args["value"] = arg.Value
	args["desc"] = arg.Desc
	return d.DB.Table(model.SiteConfigTable).Where("id=?", arg.ID).Update(args).Error
}
