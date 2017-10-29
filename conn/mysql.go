package conn

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"

	// Gorm mysql connector
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
)

type connection struct {
	DatabaseHost string `default:"localhost"`
	DatabasePort string `default:"3306"`
	DatabaseName string `default:"melhorcompra"`
	DatabaseUser string `default:"root"`
	DatabasePass string `default:"root"`
}

// InitDb init db connection
func InitDb() *gorm.DB {

	var c connection
	err := envconfig.Process("melhorcompra", &c)
	if err != nil {
		log.Fatal(err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", c.DatabaseUser, c.DatabasePass, c.DatabaseHost, c.DatabasePort, c.DatabaseName)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	/* if !db.HasTable(&models.Produto{}) {
		db.CreateTable(&models.Produto{})
	}
	if !db.HasTable(&models.Fornecedor{}) {
		db.CreateTable(&models.Fornecedor{})
	}
	if !db.HasTable(&models.Orcamento{}) {
		db.CreateTable(&models.Orcamento{})
	}
	if !db.HasTable(&models.OrcamentoFornecedor{}) {
		db.CreateTable(&models.OrcamentoFornecedor{})
	}
	if !db.HasTable(&models.OrcamentoProduto{}) {
		db.CreateTable(&models.OrcamentoProduto{})
	}
	if !db.HasTable(&models.Proposta{}) {
		db.CreateTable(&models.Proposta{})
	}
	if !db.HasTable(&models.PropostaProduto{}) {
		db.CreateTable(&models.PropostaProduto{})
	}
	if !db.HasTable(&models.Pedido{}) {
		db.CreateTable(&models.Pedido{})
	}
	if !db.HasTable(&models.PedidoProduto{}) {
		db.CreateTable(&models.PedidoProduto{})
	} */

	return db
}
