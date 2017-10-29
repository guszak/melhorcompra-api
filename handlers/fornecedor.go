package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/guszak/bestorder/models"
	"github.com/guszak/bestorder/services"
)

// ListarFornecedores list Company
func ListarFornecedores(c *gin.Context) {

	var query models.Query
	err := c.Bind(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companies, count, total, err := services.ListarFornecedores(query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else if query.Count == 1 {
		c.JSON(http.StatusOK, count)
	} else {
		c.Writer.Header().Set("X-Total-Count", strconv.Itoa(int(total)))
		c.JSON(http.StatusOK, companies)
	}
}

// VerFornecedore show Company
func VerFornecedor(c *gin.Context) {

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p, err := services.VerFornecedor(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, p)
	}
}
