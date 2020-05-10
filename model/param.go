package model

type Page struct {
	Pn    int `json:"pn" form:"pn"`
	Ps    int `json:"ps" form:"ps"`
	Total int `json:"total"`
}

type KV struct {
	K string `json:"k"`
	V string `json:"v"`
}

type Date struct {
	StartDate string `json:"start_date"`
	StartTime string `json:"start_time"`
	EndDate   string `json:"end_date"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`
	TimeZone  string `json:"time_zone"`
}
