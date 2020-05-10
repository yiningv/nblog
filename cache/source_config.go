package cache

import (
	"github.com/yiningv/nblog/model"
)

// SourceConfig cache.
var SourceConfig = &sourceConfigCache{
	nameHolder: make(map[string]*model.SourceConfig),
}

type sourceConfigCache struct {
	nameHolder map[string]*model.SourceConfig
}

func (cache *sourceConfigCache) Get(name string) *model.SourceConfig {
	return cache.nameHolder[name]
}

func (cache *sourceConfigCache) Update(newData []*model.SourceConfig) {
	nameHolder := make(map[string]*model.SourceConfig)
	for i := range newData {
		c := newData[i]
		nameHolder[c.Name] = c
	}
	cache.nameHolder = nameHolder
}

func (cache *sourceConfigCache) Replace(newData map[string]*model.SourceConfig) {
	cache.nameHolder = newData
}

func (cache *sourceConfigCache) GetAll() (all map[string]*model.SourceConfig) {
	ret := make(map[string]*model.SourceConfig)
	for s := range cache.nameHolder {
		ret[s] = cache.nameHolder[s]
	}
	return ret
}
