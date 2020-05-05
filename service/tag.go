package service

import "github.com/yiningv/nblog/model"

func (s *Service) GetTags(pn, ps int) (pager *model.TagPager, err error) {
	pager = &model.TagPager{}
	if pager.Items, pager.Page, err = s.dao.GetTags(pn, ps); err != nil {
		return
	}
	return
}
