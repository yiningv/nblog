package model

type SiteConfig struct {
	ID    int    `json:"id" gorm:"primary_key;column:id"`
	Name  string `json:"name" gorm:"column:name"`
	Type  string `json:"type" gorm:"column:type"` // string number bool
	Value string `json:"value" gorm:"column:value"`
	Desc  string `json:"desc" gorm:"column:desc"`
}

const SiteConfigTable = "site_config"

// TableName is used to identify table name in gorm
func (SiteConfig) TableName() string {
	return SiteConfigTable
}
