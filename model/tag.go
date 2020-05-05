package model

type Tag struct {
	ID           int64  `json:"id" gorm:"primary_key;column:id"`
	Name         string `json:"name" gorm:"column:name"`
	Slug         string `json:"slug" gorm:"column:slug"`
	ArticleCount int    `json:"article_count" gorm:"column:article_count"`
}

const TagTable = "tag"

// TableName is used to identify table name in gorm
func (Tag) TableName() string {
	return TagTable
}

type TagPager struct {
	Items []*Tag `json:"items"`
	Page  *Page  `json:"page"`
}

type TagRef struct {
	ID        int64 `json:"id" gorm:"primary_key;column:id"`
	TagID     int64 `json:"tag_id" gorm:"column:tag_id"`
	ArticleID int64 `json:"article_id" gorm:"column:article_id"`
}

const TagRefTable = "tag_ref"

// TableName is used to identify table name in gorm
func (TagRef) TableName() string {
	return TagRefTable
}
