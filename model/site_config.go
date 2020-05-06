package model

type SiteConfig struct {
	ID    int    `json:"id" gorm:"primary_key;column:id"`
	Name  string `json:"Name" gorm:"column:name"`
	Type  string `json:"Type" gorm:"column:type"` // string int bool image
	Value string `json:"Value" gorm:"column:value"`
	Desc  string `json:"Desc" gorm:"column:desc"`
}

const SiteConfigTable = "site_config"

// TableName is used to identify table name in gorm
func (SiteConfig) TableName() string {
	return SiteConfigTable
}
