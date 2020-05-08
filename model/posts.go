package model

type Posts struct {
	ID      int    `json:"id" gorm:"primary_key;column:id"`
	PageId  string `json:"page_id" gorm:"column:page_id;size:64;unique_index"`
	Title   string `json:"title" gorm:"column:title"`
	Slug    string `json:"slug" gorm:"column:slug;index:slug_idx"`
	Excerpt string `json:"excerpt" gorm:"column:excerpt"`
	Content string `json:"content" gorm:"column:content;type:text"`
	Status  int8   `json:"status" gorm:"column:status;index:status_idx"`
}

type PostsPager struct {
	Items []*Posts `json:"items"`
	Page  *Page    `json:"page"`
}

// TableName is used to identify table name in gorm
func (Posts) TableName() string {
	return "posts"
}
