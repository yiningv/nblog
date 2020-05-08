package model

const (
	CorrelationPostsCategory = iota
	CorrelationPostsTag
	CorrelationPostsArchive
)

// id1(posts_id) - id2(category_id)
// id1(posts_id) - id2(tag_id)
// id1(posts_id) - id2(archive_id)
type Correlation struct {
	ID   int `json:"id" gorm:"primary_key;column:id"`
	ID1  int `json:"id1" gorm:"column:id1"`
	ID2  int `json:"id2" gorm:"column:id2"`
	Type int `json:"type" gorm:"column:type"`
}

// TableName is used to identify table name in gorm
func (Correlation) TableName() string {
	return "correlation"
}
