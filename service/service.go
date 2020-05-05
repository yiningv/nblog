package service

import (
	"github.com/yiningv/nblog/conf"
	"github.com/yiningv/nblog/dao"
)

type Service struct {
	dao *dao.Dao
}

func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: dao.New(c),
	}
	return
}

func (s *Service) Close() (err error) {
	if s.dao != nil {
		err = s.dao.Close()
	}
	return
}
