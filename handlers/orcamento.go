package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guszak/melhorcompra-api/models"
	"github.com/guszak/melhorcompra-api/services"
)

// ListarOrcamentos lista os orcamentos disponiveis
func ListarOrcamentos(c *gin.Context) {

	var query models.Query
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, count, total, err := services.ListarOrcamentos(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if query.Count == 1 {
		c.JSON(http.StatusOK, count)
	} else {
		c.Writer.Header().Set("X-Total-Count", strconv.Itoa(int(total)))
		c.JSON(http.StatusOK, u)
	}
}

// VerOrcamento apresenta dados do orcamento
func VerOrcamento(c *gin.Context) {

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if u, err := services.VerOrcamento(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, u)
	}
}
