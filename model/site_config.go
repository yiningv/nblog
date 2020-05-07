package model

type SiteConfig struct {
	ID             int    `json:"-" gorm:"primary_key;column:id"`
	Name           string `json:"name" gorm:"column:name"`
	Type           string `json:"type" gorm:"column:type"` // string int bool image
	Value          string `json:"value" gorm:"column:value"`
	Desc           string `json:"desc" gorm:"column:desc"`
	LastEditedTime int64  `json:"last_edited_time" gorm:"column:last_edited_time"`
}

const SiteConfigTable = "site_config"

// TableName is used to identify table name in gorm
func (SiteConfig) TableName() string {
	return SiteConfigTable
}
