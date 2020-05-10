package cache

import (
	"github.com/yiningv/nblog/model"
	"sort"
)

// Archive cache.
var Archive = &archiveCache{
	slugHolder: make(map[string]*model.ArchivePosts),
}

type archiveCache struct {
	slugHolder map[string]*model.ArchivePosts
}

func (cache *archiveCache) Get(slug string) *model.ArchivePosts {
	ret := cache.slugHolder[slug]
	return ret
}

func (cache *archiveCache) Update(newData []*model.ArchivePosts) {
	slugHolder := make(map[string]*model.ArchivePosts)
	for i := range newData {
		archivePosts := newData[i]
		if archivePosts.SortPosts != nil {
			sort.Sort(archivePosts.SortPosts)
		}
		slugHolder[archivePosts.Archive.Slug] = archivePosts
	}
	cache.slugHolder = slugHolder
}

func (cache *archiveCache) Replace(newData map[string]*model.ArchivePosts) {
	for _, archivePosts := range newData {
		if archivePosts.SortPosts != nil {
			sort.Sort(archivePosts.SortPosts)
		}
	}
	cache.slugHolder = newData
}

func (cache *archiveCache) GetAll() (all map[string]*model.ArchivePosts) {
	ret := make(map[string]*model.ArchivePosts)
	for s := range cache.slugHolder {
		ret[s] = cache.slugHolder[s]
	}
	return ret
}
