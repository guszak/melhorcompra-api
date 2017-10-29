package models

import "time"

// OrcamentoFornecedor Model
type OrcamentoFornecedor struct {
	ID           uint64     `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	OrcamentoID  uint64     `form:"orcamento_id" json:"-"`
	FornecedorID uint64     `form:"fornecedor_id" json:"fornecedor_id"`
	CreatedAt    time.Time  `form:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `form:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `form:"deleted_at" json:"-"`
}
