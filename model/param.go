package model

type Page struct {
	Pn    int `json:"pn" form:"pn"`
	Ps    int `json:"ps" form:"ps"`
	Total int `json:"total"`
}
