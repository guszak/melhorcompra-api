package models

// Query Model
type Query struct {
	Fields string `form:"fields" json:"fields"`
	Offset int    `form:"offset" json:"offset"`
	Limit  int    `form:"limit" json:"limit"`
	Filter string `form:"filter" json:"filter"`
	Sort   string `form:"sort" json:"sort"`
	Count  int    `form:"count" json:"count"`
	Q      string `form:"q" json:"q"`
}
