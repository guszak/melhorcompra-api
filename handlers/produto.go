package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guszak/bestorder/models"
	"github.com/guszak/bestorder/services"
)

// ListarProdutos lista os produtos
func ListarProdutos(c *gin.Context) {

	var query models.Query
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, count, total, err := services.ListarProdutos(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else if query.Count == 1 {
		c.JSON(http.StatusOK, count)
	} else {
		c.Writer.Header().Set("X-Total-Count", strconv.Itoa(int(total)))
		c.JSON(http.StatusOK, p)
	}
}

// VerProduto apresenta dados do produto
func VerProduto(c *gin.Context) {

	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if p, err := services.VerProduto(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, p)
	}
}
