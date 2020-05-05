package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/pub/db"
	"sync"
)

var once sync.Once

type Dao struct {
	DB *gorm.DB
}

func New(c *conf.Config) (d *Dao) {
	once.Do(func() {
		d = &Dao{
			DB: db.NewDB(c.ORM),
		}
		d.initORM()
	})
	return
}

func (d *Dao) initORM() {
	d.DB.LogMode(true)
}

// Close close the resource.
func (d *Dao) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Ping dao ping
func (d *Dao) Ping(c context.Context) (err error) {
	if d.DB != nil {
		err = d.DB.DB().PingContext(c)
	}
	return
}
