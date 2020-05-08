package model

type Category struct {
	ID         int    `json:"id" gorm:"primary_key;column:id"`
	Name       string `json:"name" gorm:"column:name"`
	Slug       string `json:"slug" gorm:"column:slug"`
	PostsCount int    `json:"posts_count" gorm:"column:posts_count"`
}

// TableName is used to identify table name in gorm
func (Category) TableName() string {
	return "category"
}
