package model

type SourceConfig struct {
	ID             int    `json:"-" gorm:"primary_key;column:id"`
	Name           string `json:"name" gorm:"column:name;size:32"`
	Slug           string `json:"slug" gorm:"column:slug;size:64"`
	Type           string `json:"type" gorm:"column:type;size:32"`
	Table          string `json:"table" gorm:"column:table;size:64"`
	OrderNum       int    `json:"order_num" gorm:"column:order_num;size:2"`
	Desc           string `json:"desc" gorm:"column:desc;size:255"`
	LastEditedTime int64  `json:"last_edited_time" gorm:"column:last_edited_time"`
}

// TableName is used to identify table name in gorm
func (SourceConfig) TableName() string {
	return "source_config"
}
