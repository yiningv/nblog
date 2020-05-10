package cache

import (
	"github.com/yiningv/nblog/model"
	"sort"
)

// Tag cache.
var Tag = &tagCache{
	slugHolder: make(map[string]*model.TagPosts),
}

type tagCache struct {
	slugHolder map[string]*model.TagPosts
}

func (cache *tagCache) Get(slug string) *model.TagPosts {
	ret := cache.slugHolder[slug]
	return ret
}

func (cache *tagCache) Update(newData []*model.TagPosts) {
	slugHolder := make(map[string]*model.TagPosts)
	for i := range newData {
		tagPosts := newData[i]
		if tagPosts.SortPosts != nil {
			sort.Sort(tagPosts.SortPosts)
		}
		slugHolder[tagPosts.Tag.Slug] = tagPosts
	}
	cache.slugHolder = slugHolder
}

func (cache *tagCache) Replace(newData map[string]*model.TagPosts) {
	for _, tagPosts := range newData {
		if tagPosts.SortPosts != nil {
			sort.Sort(tagPosts.SortPosts)
		}
	}
	cache.slugHolder = newData
}

func (cache *tagCache) GetAll() (all map[string]*model.TagPosts) {
	ret := make(map[string]*model.TagPosts)
	for s := range cache.slugHolder {
		ret[s] = cache.slugHolder[s]
	}
	return ret
}
