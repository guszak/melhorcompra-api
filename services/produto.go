package services

import (
	"github.com/guszak/bestorder/conn"
	"github.com/guszak/bestorder/models"
)

// ListarProdutos lista os produtos
func ListarProdutos(q models.Query) ([]*models.Produto, uint64, uint64, error) {

	db := conn.InitDb()
	defer db.Close()

	var p []*models.Produto
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

// VerProduto apresenta dados do produto
func VerProduto(id uint64) (models.Produto, error) {

	db := conn.InitDb()
	defer db.Close()

	var p models.Produto

	if err := db.First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}
