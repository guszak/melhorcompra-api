package models

import "time"

// PropostaProduto Model
type PropostaProduto struct {
	ID                 uint64     `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	FornecedorID       uint64     `form:"fornecedor_id" json:"-"`
	PropostaID         uint64     `form:"proposta_id" json:"-"`
	OrcamentoProdutoID uint64     `form:"orcamento_produto_id" json:"orcamento_produto_id"`
	ClienteProdutoID   uint64     `form:"cliente_produto_id" json:"cliente_produto_id"`
	Quantidade         uint64     `form:"quantidade" json:"quantidade"`
	Preco              float32    `form:"preco" json:"preco"`
	PrazoEntrega       int        `form:"prazo_entrega" json:"prazo_entrega"`
	CreatedAt          time.Time  `form:"created_at" json:"created_at"`
	UpdatedAt          time.Time  `form:"updated_at" json:"updated_at"`
	DeletedAt          *time.Time `form:"deleted_at" json:"-"`
}
