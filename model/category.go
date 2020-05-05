package model

type Category struct {
	ID           int64  `json:"id" gorm:"primary_key;column:id"`
	Name         string `json:"name" gorm:"column:name"`
	Slug         string `json:"slug" gorm:"column:slug"`
	Desc         string `json:"desc" gorm:"column:desc"`
	ArticleCount int    `json:"article_count" gorm:"column:article_count"`
}

const CategoryTable = "category"

// TableName is used to identify table name in gorm
func (Category) TableName() string {
	return CategoryTable
}

type CategoryRef struct {
	ID         int64 `json:"id" gorm:"primary_key;column:id"`
	CategoryID int64 `json:"category_id" gorm:"column:category_id"`
	ArticleID  int64 `json:"article_id" gorm:"column:article_id"`
}

const CategoryRefTable = "category_ref"

// TableName is used to identify table name in gorm
func (CategoryRef) TableName() string {
	return CategoryRefTable
}
