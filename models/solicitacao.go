package models

// Solicitacao Model
type Solicitacao struct {
	Preco       float32   `form:"preco" json:"preco"`
	Prazo       float32   `form:"prazo" json:"prazo"`
	Negociacoes float32   `form:"negociacoes" json:"negociacoes"`
	Atrasadas   float32   `form:"atrasadas" json:"atrasadas"`
	Tempo       float32   `form:"tempo" json:"tempo"`
	Individuos  int       `form:"individuos" json:"individuos"`
	Geracoes    int       `form:"geracoes" json:"geracoes"`
	Torneio     int       `form:"torneio" json:"torneio"`
	Cruzamento  float64   `form:"cruzamento" json:"cruzamento"`
	Mutacao     float64   `form:"mutacao" json:"mutacao"`
	Orcamento   Orcamento `form:"orcamento" json:"orcamento"`
}

type Resposta struct {
	Labels    []string  `form:"labels" json:"labels"`
	Scores    []uint64  `form:"scores" json:"scores"`
	Individuo Individuo `form:"individuo" json:"individuo"`
}
