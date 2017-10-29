package services

import (
	"github.com/guszak/bestorder/conn"
	"github.com/guszak/bestorder/models"
)

// ListarFornecedores lista os fornecedores
func ListarFornecedores(q models.Query) ([]*models.Fornecedor, uint64, uint64, error) {
	db := conn.InitDb()
	defer db.Close()

	var p []*models.Fornecedor
	var count uint64
	var total uint64

	if q.Limit == 0 {
		q.Limit = 10
	}

	if q.Fields == "" {
		q.Fields = "*"
	}

	err := db.Select(q.Fields).
		Order(q.Sort).
		Offset(q.Offset).
		Limit(q.Limit).
		Find(&p).
		Count(&count).Error
	if err != nil {
		return nil, 0, 0, err
	}
	return p, count, total, nil
}

// VerFornecedor apresenta dados do fornecedor
func VerFornecedor(id uint64) (models.Fornecedor, error) {
	db := conn.InitDb()
	defer db.Close()

	var p models.Fornecedor

	if err := db.First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}
