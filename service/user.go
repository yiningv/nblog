package service

import (
	"github.com/yiningv/nblog/model"
	"github.com/yiningv/nblog/pub/util"
)

func (s *Service) GetUser(arg *model.User) (user *model.User, err error) {
	arg.Password = util.Sha1(arg.Password)

	return s.dao.GetUser(arg.Username, arg.Password)
}
