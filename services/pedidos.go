package services

import (
	"github.com/guszak/melhorcompra-api/conn"
	"github.com/guszak/melhorcompra-api/models"
)

// ListarPedidos listar Pedidos
func ListarPedidos(q models.Query) ([]*models.Pedido, uint64, uint64, error) {

	db := conn.InitDb()
	defer db.Close()

	var p []*models.Pedido
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

// VerPedido ver dados da Pedido
func VerPedido(id uint64) (models.Pedido, error) {

	db := conn.InitDb()
	defer db.Close()

	var p models.Pedido

	if err := db.First(&p, id).Error; err != nil {
		return p, err
	}
	return p, nil
}
