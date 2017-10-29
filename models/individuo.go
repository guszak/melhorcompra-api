package models

// Individuo Model
type Individuo struct {
	ID      uint64 `form:"id" json:"id"`
	Geracao uint64 `form:"geracao" json:"geracao"`
	Score   uint64 `form:"score" json:"score"`
	Genes   []Gene `form:"genes" json:"genes"`
}
