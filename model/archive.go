package model

type Archive struct {
	ID         int  `json:"id" gorm:"primary_key;column:id"`
	Year       int  `json:"year" gorm:"column:year"`
	Month      int8 `json:"month" gorm:"column:month"`
	PostsCount int  `json:"posts_count" gorm:"column:posts_count"`
}

// TableName is used to identify table name in gorm
func (Archive) TableName() string {
	return "archive"
}
