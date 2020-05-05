package model

import "time"

type Article struct {
	ID      int       `json:"id" gorm:"primary_key;column:id"`
	PageID  string    `json:"page_id" gorm:"column:page_id"`
	Title   string    `json:"title" gorm:"column:title" form:"title"`
	Slug    string    `json:"slug" gorm:"column:slug" sql:"index"`
	Excerpt string    `json:"excerpt" gorm:"column:excerpt"`
	Content string    `json:"content" gorm:"column:content"`
	Status  int8      `json:"status" gorm:"column:status"`
	Hits    int       `json:"hits" gorm:"column:hits"`
	CTime   time.Time `json:"ctime" gorm:"column:ctime"`
	MTime   time.Time `json:"mtime" gorm:"column:mtime"`
}

type ArticlePager struct {
	Items []*Article `json:"items"`
	Page  *Page      `json:"page"`
}

const ArticleTable = "article"

// TableName is used to identify table name in gorm
func (Article) TableName() string {
	return ArticleTable
}
