package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/cache"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/notion"
	"github.com/yiningv/nblog/pub/db"
	"github.com/yiningv/nblog/pub/log"
	"sync"
)

var once sync.Once

type Service struct {
	conf *conf.Config
	dao  *gorm.DB
}

func New(c *conf.Config) (s *Service) {
	s = &Service{}
	once.Do(func() {
		s.dao = db.NewDB(c.ORM)
		s.dao.LogMode(true)
	})
	return
}

func (srv *Service) Close() (err error) {
	if srv.dao != nil {
		err = srv.dao.Close()
	}
	return
}

func (srv *Service) SyncData() {
	srv.SyncSiteConfig()
	srv.SyncSourceConfig()
}

func (srv *Service) SyncSiteConfig() {
	siteConfigMap, err := notion.GetSiteConfig()
	if err != nil {
		log.Error(fmt.Sprintf("notion.GetSiteConfig failed: %v", err))
		panic(err)
	}

	siteConfigsDB, err := srv.GetSiteConfig()
	if err != nil {
		log.Error(fmt.Sprintf("srv.GetSiteConfig failed: %v", err))
		panic(err)
	}

	err = srv.BatchUpdateSiteConfig(siteConfigsDB, siteConfigMap)
	if err != nil {
		panic(err)
	}
	for k := range siteConfigMap {
		site := siteConfigMap[k]
		cache.SiteConfig.Put(site)
	}
}

func (srv *Service) SyncSourceConfig() {
	sourceConfigMap, err := notion.GetSourceConfig()
	if err != nil {
		log.Error(fmt.Sprintf("notion.GetSourceConfig failed: %v", err))
		panic(err)
	}

	sourceConfigDB, err := srv.GetSourceConfig()
	if err != nil {
		log.Error(fmt.Sprintf("srv.GetSourceConfig failed: %v", err))
		panic(err)
	}

	err = srv.BatchUpdateSourceConfig(sourceConfigDB, sourceConfigMap)
	if err != nil {
		panic(err)
	}
	for k := range sourceConfigMap {
		source := sourceConfigMap[k]
		cache.SourceConfig.Put(source)
	}
}
