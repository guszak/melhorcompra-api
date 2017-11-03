package models

import "time"

// Pedido Model
type Pedido struct {
	ID              uint64          `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	FornecedorID    uint64          `form:"fornecedor_id" json:"fornecedor_id"`
	OrcamentoID     uint64          `form:"orcamento_id" json:"orcamento_id"`
	PropostaID      uint64          `form:"proposta_id" json:"proposta_id"`
	EntregaAtrasada bool            `form:"entrega_atrasada" json:"entrega_atrasada"`
	Produtos        []PedidoProduto `form:"produtos" json:"produtos"`
	CreatedAt       time.Time       `form:"created_at" json:"created_at"`
	UpdatedAt       time.Time       `form:"updated_at" json:"updated_at"`
	DeletedAt       *time.Time      `form:"deleted_at" json:"-"`
}
