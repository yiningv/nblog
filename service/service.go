package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/pub/db"
	"github.com/yiningv/nblog/pub/log"
	"sync"
)

var (
	once sync.Once
	dao  *gorm.DB
)

func Init(c *conf.Config) {
	once.Do(func() {
		dao = db.NewDB(c.ORM)
		dao.LogMode(true)
	})
	return
}

func Close() (err error) {
	if dao != nil {
		err = dao.Close()
	}
	return
}

// 加载必要的缓存
func LoadCache() {
	err := SiteConfig.loadSiteCache()
	if err != nil {
		log.Error(fmt.Sprintf("SiteConfig.loadSiteCache failed: %v", err))
	}
	err = SourceConfig.loadSourceCache()
	if err != nil {
		log.Error(fmt.Sprintf("SourceConfig.loadSourceCache failed: %v", err))
	}
	err = Posts.loadPostsCache()
	if err != nil {
		log.Error(fmt.Sprintf("Posts.loadPostsCache failed: %v", err))
	}
	err = Tag.loadTagCache()
	if err != nil {
		log.Error(fmt.Sprintf("Tag.loadTagCache failed: %v", err))
	}
	err = Archive.loadArchiveCache()
	if err != nil {
		log.Error(fmt.Sprintf("Archive.loadArchiveCache failed: %v", err))
	}
}

// 从Notion上同步数据
func SyncData() {
	err := SiteConfig.syncSiteConfig()
	if err != nil {
		log.Error(fmt.Sprintf("SiteConfig.syncSiteConfig failed: %v", err))
	}
	err = SourceConfig.syncSourceConfig()
	if err != nil {
		log.Error(fmt.Sprintf("SourceConfig.syncSourceConfig failed: %v", err))
	}
	err = Posts.syncPosts()
	if err != nil {
		log.Error(fmt.Sprintf("Posts.syncPosts failed: %v", err))
	}
	err = Tag.syncTag()
	if err != nil {
		log.Error(fmt.Sprintf("Tag.syncTag failed: %v", err))
	}
	err = Archive.syncArchive()
	if err != nil {
		log.Error(fmt.Sprintf("Archive.syncArchive failed: %v", err))
	}
}
