package cache

import (
	"fmt"
	"github.com/bluele/gcache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"time"
)

// SourceConfig cache.
var SourceConfig = &SiteConfigCache{
	nameHolder: gcache.New(1024 * 10).LRU().Expiration(30 * time.Minute).Build(),
}

type SourceConfigCache struct {
	nameHolder gcache.Cache
}

func (cache *SourceConfigCache) Put(source *model.SourceConfig) {
	if err := cache.nameHolder.Set(source.Name, source); nil != err {
		log.Error(fmt.Sprintf("put SourceConfig [name=%d] into cache failed: %v", source.Name, err))
	}
}

func (cache *SourceConfigCache) Get(name string) *model.SourceConfig {
	ret, err := cache.nameHolder.Get(name)
	if nil != err && gcache.KeyNotFoundError != err {
		log.Error(fmt.Sprintf("get SourceConfig [name=%s] from cache failed: %v", name, err))
		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.SourceConfig)
}
