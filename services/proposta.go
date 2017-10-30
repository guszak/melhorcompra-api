package services

import (
	"github.com/guszak/melhorcompra-api/conn"
	"github.com/guszak/melhorcompra-api/models"
)

// ListarPropostas listar propostas
func ListarPropostas(q models.Query) ([]*models.Proposta, uint64, uint64, error) {

	db := conn.InitDb()
	defer db.Close()

	var p []*models.Proposta
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

// VerProposta ver dados da proposta
func VerProposta(id uint64) (models.Proposta, error) {

	db := conn.InitDb()
	defer db.Close()

	var p models.Proposta

	if err := db.First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}
