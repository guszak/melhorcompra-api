package services

import (
	"github.com/guszak/melhorcompra-api/conn"
	"github.com/guszak/melhorcompra-api/models"
)

// ListarOrcamentos listar orcamentos
func ListarOrcamentos(q models.Query) ([]*models.Orcamento, uint64, uint64, error) {

	db := conn.InitDb()
	defer db.Close()

	var p []*models.Orcamento
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

// VerOrcamento apresenta dados do orcamento
func VerOrcamento(id uint64) (models.Orcamento, error) {

	db := conn.InitDb()
	defer db.Close()

	var p models.Orcamento

	if err := db.Preload("Produtos").First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}
