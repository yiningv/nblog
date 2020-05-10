package service

import (
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/cache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/notion"
)

var SiteConfig = &siteConfigService{}

type siteConfigService struct{}

// 获取站点配置
func (srv *siteConfigService) GetSiteConfig() (siteConf []*model.SiteConfig, err error) {
	err = dao.Find(&siteConf).Error
	return
}

// 添加站点配置
func (srv *siteConfigService) AddSiteConfig(arg *model.SiteConfig) (err error) {
	return dao.Create(arg).Error
}

// 删除站点配置
func (srv *siteConfigService) DeleteSiteConfig(id int) (err error) {
	return dao.Where("id=?", id).Delete(&model.SiteConfig{}).Error
}

// 批量删除站点配置
func (srv *siteConfigService) BatchDeleteSiteConfig(ids []int) (err error) {
	return dao.Where("id IN (?)", ids).Delete(&model.SiteConfig{}).Error
}

// 修改站点配置
func (srv *siteConfigService) UpdateSiteConfig(arg *model.SiteConfig) (err error) {
	return dao.Where("id=?", arg.ID).Update(arg).Error
}

// 加载站点配置缓存
func (srv *siteConfigService) loadSiteCache() error {
	configs, err := srv.GetSiteConfig()
	if err != nil {
		return err
	}
	cache.SiteConfig.Update(configs)
	return nil
}

// 同步站点配置
func (srv *siteConfigService) syncSiteConfig() (err error) {
	var siteUpdate map[string]*model.SiteConfig
	siteUpdate, err = notion.GetSiteConfig()
	if err != nil {
		return err
	}

	siteCache := cache.SiteConfig.GetAll()

	save := make([]*model.SiteConfig, 0)
	delIds := make([]int, 0)
	for name := range siteUpdate {
		sUpdate := siteUpdate[name]
		if sCache, ok := siteCache[name]; ok {
			sUpdate.ID = sCache.ID
			// 最后更新时间有变化或者顺序有变化时，需要对数据做更新
			if sUpdate.LastEditedTime != sCache.LastEditedTime {
				save = append(save, sUpdate)
			}
			delete(siteCache, name)
		} else {
			save = append(save, sUpdate)
		}
	}
	// 缓存中剩下的数据需要删除
	for _, sCache := range siteCache {
		delIds = append(delIds, sCache.ID)
	}
	if len(save) == 0 && len(delIds) == 0 {
		return
	}
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?)", delIds).Delete(&model.SiteConfig{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			config := save[i]
			if err = tx.Save(config).Error; err != nil {
				return
			}
		}
		return
	})
	cache.SiteConfig.Replace(siteUpdate)
	return
}
