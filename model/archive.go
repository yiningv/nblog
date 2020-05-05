package model

type Archive struct {
	ID           int64 `json:"id" gorm:"primary_key;column:id"`
	Year         int   `json:"year" gorm:"column:year"`
	Month        int8  `json:"month" gorm:"column:month"`
	ArticleCount int   `json:"article_count" gorm:"column:article_count"`
}

const ArchiveTable = "archive"

// TableName is used to identify table name in gorm
func (Archive) TableName() string {
	return ArchiveTable
}

type ArchiveRef struct {
	ID        int64 `json:"id" gorm:"primary_key;column:id"`
	ArchiveID int64 `json:"archive_id" gorm:"column:archive_id"`
	ArticleID int64 `json:"article_id" gorm:"column:article_id"`
}

const ArchiveRefTable = "archive_ref"

// TableName is used to identify table name in gorm
func (ArchiveRef) TableName() string {
	return ArchiveRefTable
}
