package service

import (
	"github.com/gosimple/slug"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/cache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/notion"
)

var SourceConfig = &sourceConfigService{}

type sourceConfigService struct{}

// 获取资源配置
func (srv *sourceConfigService) GetSourceConfig() (sourceConf []*model.SourceConfig, err error) {
	err = dao.Find(&sourceConf).Error
	return
}

// 添加资源配置
func (srv *sourceConfigService) AddSourceConfig(arg *model.SourceConfig) (err error) {
	return dao.Create(arg).Error
}

// 删除资源配置
func (srv *sourceConfigService) DeleteSourceConfig(id int) (err error) {
	return dao.Where("id=?", id).Delete(&model.SourceConfig{}).Error
}

// 批量删除资源配置
func (srv *sourceConfigService) BatchDeleteSourceConfig(ids []int) (err error) {
	return dao.Where("id IN (?)", ids).Delete(&model.SourceConfig{}).Error
}

// 修改资源配置
func (srv *sourceConfigService) UpdateSourceConfig(arg *model.SourceConfig) (err error) {
	arg.Slug = slug.Make(arg.Name)
	return dao.Where("id=?", arg.ID).Update(arg).Error
}

// 加载资源配置缓存
func (srv *sourceConfigService) loadSourceCache() error {
	configs, err := srv.GetSourceConfig()
	if err != nil {
		return err
	}
	cache.SourceConfig.Update(configs)
	return nil
}

// 同步资源配置
func (srv *sourceConfigService) syncSourceConfig() (err error) {
	var sourceUpdate map[string]*model.SourceConfig
	sourceUpdate, err = notion.GetSourceConfig()
	if err != nil {
		return err
	}

	sourceCache := cache.SourceConfig.GetAll()

	save := make([]*model.SourceConfig, 0)
	delIds := make([]int, 0)
	for name := range sourceUpdate {
		sUpdate := sourceUpdate[name]
		if sCache, ok := sourceCache[name]; ok {
			sUpdate.ID = sCache.ID
			// 最后更新时间有变化或者顺序有变化时，需要对数据做更新
			if sUpdate.LastEditedTime != sCache.LastEditedTime ||
				sCache.OrderNum != sCache.OrderNum {
				save = append(save, sUpdate)
			}
			delete(sourceCache, name)
		} else {
			save = append(save, sUpdate)
		}
	}
	// 缓存中剩下的数据需要删除
	for _, sCache := range sourceCache {
		delIds = append(delIds, sCache.ID)
	}
	if len(save) == 0 && len(delIds) == 0 {
		return
	}
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delIds) > 0 {
			if err = tx.Where("id IN (?)", delIds).Delete(&model.SourceConfig{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			config := save[i]
			config.Slug = slug.Make(config.Name)
			if err = tx.Save(config).Error; err != nil {
				return
			}
		}
		return
	})
	cache.SourceConfig.Replace(sourceUpdate)
	return
}

// 修改资源配置
func (srv *sourceConfigService) BatchUpdateSourceConfig(fromDB []*model.SourceConfig, sourceConfigMap map[string]*model.SourceConfig) (err error) {
	save := make([]*model.SourceConfig, 0)
	delNames := make([]string, 0)
	for i := range fromDB {
		scDB := fromDB[i]
		if sc, ok := sourceConfigMap[scDB.Name]; ok {
			sc.ID = scDB.ID
			if sc.LastEditedTime != scDB.LastEditedTime || sc.OrderNum != scDB.OrderNum {
				save = append(save, sc)
			}
			delete(sourceConfigMap, scDB.Name)
		} else {
			delNames = append(delNames, scDB.Name)
		}
	}
	for k := range sourceConfigMap {
		save = append(save, sourceConfigMap[k])
	}
	if len(save) == 0 && len(delNames) == 0 {
		return
	}
	err = dao.Transaction(func(tx *gorm.DB) (err error) {
		if len(delNames) > 0 {
			if err = tx.Where("name IN (?)", delNames).Delete(&model.SourceConfig{}).Error; err != nil {
				return
			}
		}

		for i := range save {
			config := save[i]
			config.Slug = slug.Make(config.Name)
			if err = tx.Save(config).Error; err != nil {
				return
			}
		}
		return
	})
	return
}
