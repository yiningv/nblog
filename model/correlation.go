package model

const (
	CorrelationPostsTag = iota
	CorrelationPostsArchive
)

// id1(posts_id) - id2(tag_id)
// id1(posts_id) - id2(archive_id)
// str1(page_id) - str2(tag_slug)
// str1(page_id) - str2(archive_slug)
type Correlation struct {
	ID   int    `json:"id" gorm:"primary_key;column:id"`
	ID1  int    `json:"id1" gorm:"column:id1"`
	ID2  int    `json:"id2" gorm:"column:id2"`
	Str1 string `json:"str1" gorm:"column:str1;size:255"`
	Str2 string `json:"str2" gorm:"column:str2;size:255"`
	Type int    `json:"type" gorm:"column:type;size:2"`
}

// TableName is used to identify table name in gorm
func (Correlation) TableName() string {
	return "correlation"
}
