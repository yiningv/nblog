package cache

import (
	"github.com/yiningv/nblog/model"
	"sort"
)

// Posts cache.
var Posts = &postsCache{
	data:       make([]*model.Posts, 0),
	slugHolder: make(map[string]*model.Posts),
}

type postsCache struct {
	data       model.SortPosts // 排序用
	slugHolder map[string]*model.Posts
}

func (cache *postsCache) Get(pageId string) *model.Posts {
	return cache.slugHolder[pageId]
}

func (cache *postsCache) Update(save []*model.Posts) {
	var sortData model.SortPosts
	sortData = save
	if sortData != nil {
		sort.Sort(sortData)
	}

	pageIdHolder := make(map[string]*model.Posts)
	for i := range save {
		c := save[i]
		pageIdHolder[c.Slug] = c
	}
	cache.slugHolder = pageIdHolder
	cache.data = sortData

}

func (cache *postsCache) GetAll() map[string]*model.Posts {
	ret := make(map[string]*model.Posts)
	for s := range cache.slugHolder {
		ret[s] = cache.slugHolder[s]
	}
	return ret
}
