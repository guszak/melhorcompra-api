package models

import "time"

// PedidoProduto Model
type PedidoProduto struct {
	ID         uint64     `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	PedidoID   uint64     `form:"pedido_id" json:"-"`
	ProdutoID  uint64     `form:"produto_id" json:"-"`
	Quantidade uint64     `form:"quantidade" json:"quantidade"`
	Preco      float32    `form:"preco" json:"preco"`
	CreatedAt  time.Time  `form:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `form:"updated_at" json:"updated_at"`
	DeletedAt  *time.Time `form:"deleted_at" json:"-"`
}
