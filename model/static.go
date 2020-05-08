package model

// 静态类数据，直接保存json的字符串
type Static struct {
	ID             int    `json:"-" gorm:"primary_key;column:id"`
	Type           string `json:"type" gorm:"column:type;size:32;index:type_idx"`
	PageId         string `json:"page_id" gorm:"column:page_id;size:64;unique_index"`
	JSONData       string `json:"json_data" gorm:"column:json_data;type:text"`
	LastEditedTime int64  `json:"last_edited_time" gorm:"column:last_edited_time"`
}

// TableName is used to identify table name in gorm
func (Static) TableName() string {
	return "static"
}
