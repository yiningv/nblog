package model

type Archive struct {
	ID         int    `json:"id" gorm:"primary_key;column:id"`
	Slug       string `json:"slug" gorm:"column:slug;size:255"`
	Year       string `json:"year" gorm:"column:year;size:4"`
	Month      string `json:"month" gorm:"column:month;size:2"`
	PostsCount int    `json:"posts_count" gorm:"column:posts_count"`
}

// TableName is used to identify table name in gorm
func (Archive) TableName() string {
	return "archive"
}

type ArchivePager struct {
	Items []*Archive `json:"items"`
	Page  *Page      `json:"page"`
}

type ArchivePosts struct {
	Archive   *Archive
	SortPosts SortPosts
}
