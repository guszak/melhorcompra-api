package models

import "time"

// Fornecedor Model
type Fornecedor struct {
	ID        uint64     `gorm:"primary_key;AUTO_INCREMENT" form:"id" json:"id"`
	Nome      string     `form:"nome" json:"nome"`
	CreatedAt time.Time  `form:"created_at" json:"created_at"`
	UpdatedAt time.Time  `form:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `form:"deleted_at" json:"-"`
}
