package models

import "time"

// OrcamentoProduto Model
type OrcamentoProduto struct {
	ID          uint64     `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	OrcamentoID uint64     `form:"orcamento_id" json:"-"`
	ProdutoID   uint64     `form:"produto_id" json:"produto_id"`
	Quantidade  uint64     `form:"quantidade" json:"quantidade"`
	CreatedAt   time.Time  `form:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `form:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time `form:"deleted_at" json:"-"`
}
