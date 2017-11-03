package models

import "time"

// Produto Model
type Produto struct {
	ID           uint64     `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	FornecedorID uint64     `form:"fornecedor_id" json:"fornecedor_id"`
	Descricao    string     `form:"descricao" json:"descricao"`
	Unidade      string     `form:"unidade" json:"unidade"`
	CreatedAt    time.Time  `form:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `form:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `form:"deleted_at" json:"-"`
}
