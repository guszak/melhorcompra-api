package main

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guszak/bestorder/conn"
	"github.com/guszak/bestorder/handlers"
	"github.com/guszak/bestorder/models"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	port := "3001"
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}

	g := gin.Default()
	g.Use(cors.Default())
	g.GET("/popular", gerarPopulacaoInicialAleatoria)
	//g.POST("/orcamento", obterMelhorCompra)
	g.GET("/produtos", handlers.ListarProdutos)
	g.GET("/produtos/:id", handlers.VerProduto)
	g.GET("/fornecedores", handlers.ListarFornecedores)
	g.GET("/fornecedores/:id", handlers.VerFornecedor)
	g.GET("/orcamentos", handlers.ListarOrcamentos)
	g.GET("/orcamentos/:id", handlers.VerOrcamento)
	g.GET("/propostas", handlers.ListarPropostas)
	g.GET("/propostas/:id", handlers.VerProposta)
	g.GET("/pedidos", handlers.ListarPedidos)
	g.GET("/pedidos/:id", handlers.VerPedido)
	http.Handle("/", g)
	http.ListenAndServe(":"+port, nil)
}

// GetAuth user account
func gerarPopulacaoInicialAleatoria(c *gin.Context) {

	db := conn.InitDb()
	defer db.Close()

	// Fornecedores
	for i := 0; i < 10; i++ {
		//var fornecedor =
		db.Create(models.Fornecedor{
			Nome: "Fornecedor " + strconv.Itoa(i)})
	}

	// Produtos
	for i := 0; i < 20; i++ {
		//var produto = models.Produto{Descricao: "Produto " + strconv.Itoa(i), Unidade: "Unidade"}
		db.Create(&models.Produto{
			Descricao: "Produto " + strconv.Itoa(i),
			Unidade:   "Unidade"})
	}

	// Orcamentos
	var produtos []*models.Produto
	var fornecedores []*models.Fornecedor
	db.Find(&produtos)
	db.Find(&fornecedores)
	for i := 0; i < 10; i++ {
		var orcamento = models.Orcamento{Descricao: "Teste"}
		for j := range produtos {
			var p = models.OrcamentoProduto{ProdutoID: produtos[j].ID, Quantidade: uint64(rand.Intn(100))}
			orcamento.Produtos = append(orcamento.Produtos, p)
		}
		for j := range fornecedores {
			var f = models.OrcamentoFornecedor{FornecedorID: fornecedores[j].ID}
			orcamento.Fornecedores = append(orcamento.Fornecedores, f)
		}
		db.Create(&orcamento)
	}

	// Propostas
	var orcamentos []*models.Orcamento
	db.Preload("Produtos").Preload("Fornecedores").Find(&orcamentos)
	for i, orcamento := range orcamentos {

		for f := range orcamento.Fornecedores {

			var proposta = models.Proposta{
				Descricao:    "Teste",
				FornecedorID: orcamento.Fornecedores[f].FornecedorID,
				OrcamentoID:  orcamento.ID}
			for j := range orcamentos[i].Produtos {
				var p = models.PropostaProduto{
					FornecedorID:       orcamento.Fornecedores[f].FornecedorID,
					OrcamentoProdutoID: orcamento.ID,
					ClienteProdutoID:   orcamento.Produtos[j].ProdutoID,
					Quantidade:         orcamento.Produtos[j].Quantidade,
					Preco:              rand.Float32(),
					PrazoEntrega:       rand.Int()}
				proposta.Produtos = append(proposta.Produtos, p)
			}
			db.Create(&proposta)
		}
	}

	//var pedido = Pedido{Descricao: "Teste"}
	//db.Create(&pedido)
}

// Apresenta ao comprador a predição da melhor compra com base no histórico de compras
func obterMelhorCompra(c *gin.Context) {

	// Preenche struct do orçamento com o conteudo do body da requisição
	var orcamento models.Orcamento
	c.Bind(&orcamento)

	// Com base nos itens do orçamento, monta a população inicial com score para cada produto
	var individuos []models.Individuo
	obterPopulacaoInicial(orcamento, &individuos)

	// Obtem o score da populacao inicial com base negociacoes do fornecedor
	gerarScorePopulacaoInicial(&individuos)

	// Combina os individuos em busca de atingir uma nova população
	//combinarIndividuos(&individuos)
}

// Com base nos itens do orçamento, monta a população inicial com score para cada produto
func obterPopulacaoInicial(orcamento models.Orcamento, individuos *[]models.Individuo) {

	// Conexão com o banco de dados
	db := conn.InitDb()
	defer db.Close()

	// Obtem o id de todos os produtos solicitados no orçamento e monta uma lista
	var produtosOrcamento []uint64
	for i := range orcamento.Produtos {
		produtosOrcamento = append(produtosOrcamento, orcamento.Produtos[i].ProdutoID)
	}

	// Busca todos os orçamentos já realizados com os produtos solicitados
	var orcamentos []uint64
	db.
		Model(&models.OrcamentoProduto{}).
		Where("produto_id in (?)", produtosOrcamento).
		Group("orcamento_id").
		Pluck("orcamento_id", &orcamentos)

	// Busca todos os fornecedores associados a esses orçamentos e que vão compor a população inicial
	var fornecedores []uint64
	db.
		Model(&models.OrcamentoFornecedor{}).
		Where("orcamento_id in (?)", orcamentos).
		Group("fornecedor_id").
		Pluck("fornecedor_id", &fornecedores)

	// Obtem o ultimo produto proposto para cada produto solicitado no orçamento
	var produtos []*models.PropostaProduto
	db.
		Where("cliente_produto_id in (?)", produtosOrcamento).
		Group("cliente_produto_id").
		Order("fornecedor_id desc, created_at desc").
		Find(&produtos)

	var gene models.Gene
	var found bool
	var fornecedorAtual uint64
	var individuoID uint64
	// Inviduos
	for x := range produtos {

		if produtos[x].FornecedorID == fornecedorAtual {
			continue
		}

		fornecedorAtual = produtos[x].FornecedorID

		var individuo = models.Individuo{
			ID:      individuoID,
			Geracao: 1}

		// Adiciona novo individuo na populacao inicial
		*individuos = append(*individuos, individuo)

		// Adiciona os genes ao individuo
		for i := range orcamento.Produtos {

			found = false

			// Busca algum historico de orçamento para o produto e associa os dados da ultima proposta
			for j := range produtos {
				if orcamento.Produtos[j].ProdutoID == produtos[j].ClienteProdutoID {
					gene = models.Gene{
						OrcamentoProdutoID: i,
						ClienteProdutoID:   orcamento.Produtos[i].ProdutoID,
						IndividuoID:        individuoID,
						FornecedorID:       fornecedorAtual,
						Preco:              produtos[j].Preco,
						PrazoEntrega:       produtos[j].PrazoEntrega}
					found = true
					break
				}
			}

			// Caso, nao encontre algum historico de cotacao, apenas associa as informacoes basicas
			if !found {
				gene = models.Gene{
					OrcamentoProdutoID: i,
					ClienteProdutoID:   orcamento.Produtos[i].ProdutoID,
					IndividuoID:        individuoID,
					FornecedorID:       fornecedorAtual}
			}

			// Adiciona o gene ao individuo
			(*individuos)[individuoID].Genes = append((*individuos)[individuoID].Genes, gene)
		}
		individuoID++
	}
}

// Obtem o score da populacao inicial com base negociacoes do fornecedor
func gerarScorePopulacaoInicial(individuos *[]models.Individuo) {

	// Conexão com o banco de dados
	//db := conn.InitDb()
	//defer db.Close()

	// Entregas Atrasadas
	// rows, _ := db.
	// 	Model(&models.Pedido{}).
	// 	Select("fornecedor_id, count(1) as total").
	// 	Where("fornecedor_id in (?)", fornecedores).
	// 	Where("entrega_atrasada = ?", true).
	// 	Group("fornecedor_id").
	// 	Rows()

	// // Busca maior e menor valor
	// for rows.Next() {
	// 	//rows.Columns.fornecedor_id
	// }

	// // Total de Pedidos
	// db.
	// 	Model(&models.Pedido{}).
	// 	Select("fornecedor_id, count(1) as total").
	// 	Where("fornecedor_id in (?)", fornecedores).
	// 	Group("fornecedor_id").
	// 	Rows()
	//for rows.Next() {
	//...
	//}

	// Tempo de negociacao
	//db.
	//Model(&Pedido{}).
	//Select("fornecedor_id, count(1) as total").
	//Where("fornecedor_id in (?)", fornecedores).
	//Group("fornecedor_id")

	// Score por produto
	for i := range (*individuos)[0].Genes {

		// Preco
		var menorPreco, maiorPreco, notaPreco float32
		var menorPrazo, maiorPrazo, notaPrazo int

		// Maior Preço
		for j := range *individuos {
			if (*individuos)[j].Genes[i].Preco > maiorPreco {
				maiorPreco = (*individuos)[j].Genes[i].Preco
			}
		}

		// Menor Preço
		menorPreco = maiorPreco
		for j := range *individuos {
			if (*individuos)[j].Genes[i].Preco < menorPreco {
				menorPreco = (*individuos)[j].Genes[i].Preco
			}
		}

		for j := range *individuos {
			notaPreco = ((*individuos)[j].Genes[i].Preco - menorPreco) * (10 / (maiorPreco - menorPreco))
			(*individuos)[j].Genes[i].Score += uint64(notaPreco)
		}

		// Maior Prazo
		for j := range *individuos {
			if (*individuos)[j].Genes[i].PrazoEntrega > maiorPrazo {
				maiorPrazo = (*individuos)[j].Genes[i].PrazoEntrega
			}
		}

		// Menor Prazo
		menorPrazo = maiorPrazo
		for j := range *individuos {
			if (*individuos)[j].Genes[i].PrazoEntrega < menorPrazo {
				menorPrazo = (*individuos)[j].Genes[i].PrazoEntrega
			}
		}

		for j := range *individuos {
			notaPrazo = ((*individuos)[j].Genes[i].PrazoEntrega - menorPrazo) * (10 / (maiorPrazo - menorPrazo))
			(*individuos)[j].Genes[i].Score += uint64(notaPrazo)
		}
	}
}
