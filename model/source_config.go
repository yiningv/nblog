package model

type SourceConfig struct {
	ID    int    `json:"id" gorm:"primary_key;column:id"`
	Name  string `json:"Name" gorm:"column:name"`
	Type  string `json:"Type" gorm:"column:type"` // string number bool
	Table string `json:"Table" gorm:"column:table"`
	Desc  string `json:"Desc" gorm:"column:desc"`
}

const SourceConfigTable = "source_config"

// TableName is used to identify table name in gorm
func (SourceConfig) TableName() string {
	return SourceConfigTable
}
