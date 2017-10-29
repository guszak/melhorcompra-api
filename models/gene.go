package models

// Gene Model
type Gene struct {
	ID                 uint64  `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	IndividuoID        uint64  `form:"inviduo_id" json:"-"`
	FornecedorID       uint64  `form:"fornecedor_id" json:"-"`
	OrcamentoProdutoID int     `form:"orcamento_produto_id" json:"orcamento_produto_id"`
	ClienteProdutoID   uint64  `form:"cliente_produto_id" json:"cliente_produto_id"`
	Quantidade         uint64  `form:"quantidade" json:"quantidade"`
	Preco              float32 `form:"preco" json:"preco"`
	PrazoEntrega       int     `form:"prazo_entrega" json:"prazo_entrega"`
	Score              uint64  `form:"score" json:"score"`
}
