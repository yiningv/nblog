package model

type SiteConfig struct {
	ID             int    `json:"-" gorm:"primary_key;column:id"`
	Name           string `json:"name" gorm:"column:name;size:32"`
	Type           string `json:"type" gorm:"column:type;size:16"` // string int bool image bool
	Value          string `json:"value" gorm:"column:value;size:255"`
	Desc           string `json:"desc" gorm:"column:desc;size:255"`
	LastEditedTime int64  `json:"last_edited_time" gorm:"column:last_edited_time"`
}

// TableName is used to identify table name in gorm
func (SiteConfig) TableName() string {
	return "site_config"
}
