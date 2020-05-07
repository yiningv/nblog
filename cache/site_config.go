package cache

import (
	"fmt"
	"github.com/bluele/gcache"
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/log"
	"time"
)

// SiteConfig cache.
var SiteConfig = &SiteConfigCache{
	nameHolder: gcache.New(1024 * 10).LRU().Expiration(30 * time.Minute).Build(),
}

type SiteConfigCache struct {
	nameHolder gcache.Cache
}

func (cache *SiteConfigCache) Put(site *model.SiteConfig) {
	if err := cache.nameHolder.Set(site.Name, site); nil != err {
		log.Error(fmt.Sprintf("put SiteConfig [name=%s] into cache failed: %v", site.Name, err))
	}
}

func (cache *SiteConfigCache) Get(name string) *model.SiteConfig {
	ret, err := cache.nameHolder.Get(name)
	if nil != err && gcache.KeyNotFoundError != err {
		log.Error(fmt.Sprintf("get SiteConfig [name=%s] from cache failed: %v", name, err))
		return nil
	}
	if nil == ret {
		return nil
	}

	return ret.(*model.SiteConfig)
}
