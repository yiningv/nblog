package model

import "time"

type Posts struct {
	ID             int       `json:"id" gorm:"primary_key;column:id"`
	PageId         string    `json:"page_id" gorm:"column:page_id;size:64;unique_index"`
	Title          string    `json:"title" gorm:"column:title;size:128"`
	Slug           string    `json:"slug" gorm:"column:slug;size:255;index:slug_idx"`
	Excerpt        string    `json:"excerpt" gorm:"column:excerpt;type:text"`
	Content        string    `json:"content" gorm:"column:content;type:text"`
	Tags           string    `json:"tags" gorm:"column:tags;type:text"`
	PublishedTime  string    `json:"published_time" gorm:"-"`
	PTime          time.Time `json:"p_time" gorm:"column:p_time"`
	Status         string    `json:"status" gorm:"column:status;size:16;index:status_idx"`
	LastEditedTime int64     `json:"last_edited_time" gorm:"column:last_edited_time"`
}

type PostsPager struct {
	Items []*Posts `json:"items"`
	Page  *Page    `json:"page"`
}

// TableName is used to identify table name in gorm
func (Posts) TableName() string {
	return "posts"
}

type SortPosts []*Posts

func (s SortPosts) Len() int           { return len(s) }
func (s SortPosts) Less(i, j int) bool { return s[i].PTime.After(s[j].PTime) }
func (s SortPosts) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
