package model

type SourceConfig struct {
	ID    int    `json:"id" gorm:"primary_key;column:id"`
	Name  string `json:"name" gorm:"column:name"`
	Type  string `json:"type" gorm:"column:type"` // string number bool
	Table string `json:"table" gorm:"column:table"`
	Desc  string `json:"desc" gorm:"column:desc"`
}

const SourceConfigTable = "source_config"

// TableName is used to identify table name in gorm
func (SourceConfig) TableName() string {
	return SourceConfigTable
}
