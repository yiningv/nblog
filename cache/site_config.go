package cache

import (
	"github.com/yiningv/nblog/model"
)

// SiteConfig cache.
var SiteConfig = &siteConfigCache{
	nameHolder: make(map[string]*model.SiteConfig),
}

type siteConfigCache struct {
	nameHolder map[string]*model.SiteConfig
}

func (cache *siteConfigCache) Get(name string) *model.SiteConfig {
	return cache.nameHolder[name]
}

func (cache *siteConfigCache) Update(newData []*model.SiteConfig) {
	nameHolder := make(map[string]*model.SiteConfig)
	for i := range newData {
		c := newData[i]
		nameHolder[c.Name] = c
	}
	cache.nameHolder = nameHolder
}

func (cache *siteConfigCache) Replace(newData map[string]*model.SiteConfig) {
	cache.nameHolder = newData
}

func (cache *siteConfigCache) GetAll() (all map[string]*model.SiteConfig) {
	ret := make(map[string]*model.SiteConfig)
	for s := range cache.nameHolder {
		ret[s] = cache.nameHolder[s]
	}
	return ret
}
