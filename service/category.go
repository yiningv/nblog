package service

import "github.com/yiningv/nblog/model"

// 获取分类列表
func (s *Service) GetCategories() (cates []*model.Category, err error) {
	if cates, err = s.dao.GetCategories(); err != nil {
		return
	}
	return
}

// 删除分类
func (s *Service) DeleteCategories(ids []int64) (err error) {
	switch len(ids) {
	case 0:
		return
	case 1:
		if err = s.dao.DeleteCategory(ids[0]); err != nil {
			return
		}
	default:
		if err = s.dao.BatchDeleteCategory(ids); err != nil {
			return
		}
	}
	return
}

// 更新分类
func (s *Service) UpdateCategory(arg *model.Category) (err error) {
	err = s.dao.UpdateCategory(arg)
	return
}
