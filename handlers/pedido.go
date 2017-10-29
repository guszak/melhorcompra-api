package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guszak/bestorder/models"
	"github.com/guszak/bestorder/services"
)

// ListarPedidos lista os pedidos
func ListarPedidos(c *gin.Context) {

	var query models.Query
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, count, total, err := services.ListarPedidos(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if query.Count == 1 {
		c.JSON(http.StatusOK, count)
	} else {
		c.Writer.Header().Set("X-Total-Count", strconv.Itoa(int(total)))
		c.JSON(http.StatusOK, u)
	}
}

// VerPedido apresenta dados do pedido
func VerPedido(c *gin.Context) {

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if u, err := services.VerPedido(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, u)
	}
}
