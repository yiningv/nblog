package model

type Tag struct {
	ID         int    `json:"id" gorm:"primary_key;column:id"`
	Name       string `json:"name" gorm:"column:name;size:128"`
	Slug       string `json:"slug" gorm:"column:slug;size:128"`
	PostsCount int    `json:"posts_count" gorm:"column:posts_count"`
}

// TableName is used to identify table name in gorm
func (Tag) TableName() string {
	return "tag"
}

type TagPager struct {
	Items []*Tag `json:"items"`
	Page  *Page  `json:"page"`
}

type TagPosts struct {
	Tag       *Tag
	SortPosts SortPosts
}
