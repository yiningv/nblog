package service

import "github.com/yiningv/nblog/model"

func (s *Service) GetArticles(pn int, ps int) (pager *model.ArticlePager, err error) {
	pager = &model.ArticlePager{}
	if pager.Items, pager.Page, err = s.dao.GetArticles(pn, ps); err != nil {
		return
	}
	return
}
