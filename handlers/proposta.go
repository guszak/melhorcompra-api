package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guszak/melhorcompra-api/models"
	"github.com/guszak/melhorcompra-api/services"
)

// ListarPropostas lista as propostas
func ListarPropostas(c *gin.Context) {

	var query models.Query
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, count, total, err := services.ListarPropostas(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if query.Count == 1 {
		c.JSON(http.StatusOK, count)
	} else {
		c.Writer.Header().Set("X-Total-Count", strconv.Itoa(int(total)))
		c.JSON(http.StatusOK, u)
	}
}

// VerProposta apresenta dados da proposta
func VerProposta(c *gin.Context) {

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if u, err := services.VerProposta(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, u)
	}
}
