package model

type SourceConfig struct {
	ID             int    `json:"-" gorm:"primary_key;column:id"`
	Name           string `json:"name" gorm:"column:name"`
	Type           string `json:"type" gorm:"column:type"` // string number bool
	Table          string `json:"table" gorm:"column:table"`
	Desc           string `json:"desc" gorm:"column:desc"`
	LastEditedTime int64  `json:"last_edited_time" gorm:"column:last_edited_time"`
}

const SourceConfigTable = "source_config"

// TableName is used to identify table name in gorm
func (SourceConfig) TableName() string {
	return SourceConfigTable
}
