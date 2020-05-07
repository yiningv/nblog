package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
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

func (srv *Service) InitData() {
	// 初始化时从数据库中读取站点配置和资源配置
	dbSiteConfs, err := srv.GetSiteConfig()
	if err != nil {
		log.Error(fmt.Sprintf("router.Init error: %v", err))
		panic(err)
	}
	if len(dbSiteConfs) == 0 {
		// 从notion上读取配置文件，更新数据库和缓存
		notion.GetSiteConfig()
	}

}
