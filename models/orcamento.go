package models

import "time"

// Orcamento Model
type Orcamento struct {
	ID           uint64                `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	Descricao    string                `form:"descricao" json:"descricao"`
	Produtos     []OrcamentoProduto    `form:"produtos" json:"produtos"`
	Fornecedores []OrcamentoFornecedor `form:"fornecedores" json:"fornecedores"`
	CreatedAt    time.Time             `form:"created_at" json:"created_at"`
	UpdatedAt    time.Time             `form:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time            `form:"deleted_at" json:"-"`
}
