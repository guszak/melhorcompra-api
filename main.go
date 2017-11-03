package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/guszak/melhorcompra-api/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/kelseyhightower/envconfig"
)

// Parâmetros de conexão com o banco de dados
type connection struct {
	DatabaseHost string `default:"localhost"`
	DatabasePort string `default:"3306"`
	DatabaseName string `default:"melhorcompra"`
	DatabaseUser string `default:"root"`
	DatabasePass string `default:"root"`
}

// Totais por Fornecedor
type Total struct {
	FornecedorID uint64  `form:"fornecedor_id" json:"fornecedor_id"`
	Atrasos      float32 `form:"atrasos" json:"atrasos"`
	Compras      float32 `form:"compras" json:"compras"`
	Tempo        float32 `form:"tempo" json:"tempo"`
}

// Função inicial da api
func main() {
	port := "3001"
	if os.Getenv("HTTP_PLATFORM_PORT") != "" {
		port = os.Getenv("HTTP_PLATFORM_PORT")
	}

	g := gin.Default()
	g.Use(cors.Default())
	g.GET("/popular", gerarPopulacaoInicialAleatoria)
	g.POST("/orcamento", obterMelhorCompra)
	g.GET("/produtos", listarProdutos)
	http.Handle("/", g)
	http.ListenAndServe(":"+port, nil)
}

// Inicializa conexão com o banco de dados
func initDb() *gorm.DB {

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
	db.LogMode(true)

	return db
}

// ListarProdutos lista os produtos
func listarProdutos(c *gin.Context) {

	db := initDb()
	defer db.Close()

	var p []*models.Produto

	err := db.Find(&p).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, p)
	}
}

// GetAuth user account
func gerarPopulacaoInicialAleatoria(c *gin.Context) {

	db := initDb()
	defer db.Close()

	db.DropTableIfExists(&models.Produto{})
	db.CreateTable(&models.Produto{})

	db.DropTableIfExists(&models.Fornecedor{})
	db.CreateTable(&models.Fornecedor{})

	db.DropTableIfExists(&models.Orcamento{})
	db.CreateTable(&models.Orcamento{})

	db.DropTableIfExists(&models.OrcamentoFornecedor{})
	db.CreateTable(&models.OrcamentoFornecedor{})

	db.DropTableIfExists(&models.OrcamentoProduto{})
	db.CreateTable(&models.OrcamentoProduto{})

	db.DropTableIfExists(&models.Proposta{})
	db.CreateTable(&models.Proposta{})

	db.DropTableIfExists(&models.PropostaProduto{})
	db.CreateTable(&models.PropostaProduto{})

	db.DropTableIfExists(&models.Pedido{})
	db.CreateTable(&models.Pedido{})

	db.DropTableIfExists(&models.PedidoProduto{})
	db.CreateTable(&models.PedidoProduto{})

	// Fornecedores
	for i := 0; i < 100; i++ {
		//var fornecedor =
		db.Create(&models.Fornecedor{
			Nome: "Fornecedor " + strconv.Itoa(i)})
	}

	// Produtos
	for i := 0; i < 50; i++ {
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

		for _, f := range orcamento.Fornecedores {

			var proposta = models.Proposta{
				Descricao:    "Teste",
				FornecedorID: f.FornecedorID,
				OrcamentoID:  orcamento.ID}
			for j := range orcamentos[i].Produtos {
				var p = models.PropostaProduto{
					FornecedorID:       f.FornecedorID,
					OrcamentoProdutoID: orcamento.ID,
					ClienteProdutoID:   orcamento.Produtos[j].ProdutoID,
					Quantidade:         orcamento.Produtos[j].Quantidade,
					Preco:              rand.Float32() * 100,
					PrazoEntrega:       rand.Intn(90)}
				proposta.Produtos = append(proposta.Produtos, p)
			}
			db.Create(&proposta)

			if rand.Float32() < 0.3 {
				var pedido = models.Pedido{
					FornecedorID:    f.FornecedorID,
					OrcamentoID:     orcamento.ID,
					EntregaAtrasada: (rand.Float32() < 0.5),
					PropostaID:      proposta.ID}
				db.Create(&pedido)
			}
		}
	}

	//var pedido = Pedido{Descricao: "Teste"}
	//db.Create(&pedido)
}

// Apresenta ao comprador a predição da melhor compra com base no histórico de compras
func obterMelhorCompra(c *gin.Context) {

	// Preenche struct do orçamento com o conteudo do body da requisição
	var solicitacao models.Solicitacao
	c.Bind(&solicitacao)
	orcamento := solicitacao.Orcamento

	// Com base nos itens do orçamento, monta a população inicial com score para cada produto
	var populacao []models.Individuo
	obterPopulacaoInicial(orcamento, &populacao)

	// Obtem o score da populacao inicial com base negociacoes do fornecedor
	gerarScorePopulacaoInicial(&populacao, solicitacao)

	// Combina os populacao em busca de atingir uma nova população
	geracoes := solicitacao.Geracoes // Numero de gerações da população
	tamTorneio := solicitacao.Torneio
	tamPop := len(populacao)            // Tamanho da população
	probCruz := solicitacao.Cruzamento  // probabilidade de cruzamento
	probMut := solicitacao.Mutacao      // probabilidade de mutação
	tamGenes := len(populacao[0].Genes) // numero de genes do individuo
	var resposta models.Resposta
	for i := 0; i < geracoes; i++ {
		for j := 0; j < tamTorneio; j++ {

			// calcula a probabilidade de cruzamento
			if rand.Float64() < probCruz {

				// escolhe dois pais
				indicePai1 := rand.Intn(tamPop)
				indicePai2 := indicePai1
				// garante que os índices dos pais não são iguais
				for indicePai1 == indicePai2 {
					indicePai2 = rand.Intn(tamPop)
				}

				var filho models.Individuo
				// aplica o cruzamento de 1 ponto
				cruzamento(indicePai1, indicePai2, populacao, &filho)

				// calcula a probabilidade de mutação
				if rand.Float64() < probMut {
					mutacao(&filho, tamGenes)
				}

				scorePai := obterPontuacao(populacao[indicePai1])
				scoreFilho := obterPontuacao(filho)

				/*
					se a pontuação (score) do filho for melhor do
					que o pai1, então substitui o pai 1 pelo filho
				*/
				if scoreFilho > scorePai {
					// faz a cópia dos genes do filho para o pai
					for k := 0; k < tamGenes; k++ {
						populacao[indicePai1].Genes[k] = filho.Genes[k]
					}
				}
			}
		}

		resposta.Labels = append(resposta.Labels, "Geracao: "+strconv.Itoa(i+1))
		fmt.Println("Geracao: " + strconv.Itoa(i+1))
		fmt.Println("Melhor: ")

		indiceMelhor := obterMelhor(populacao)
		scoreMelhor := obterPontuacao(populacao[indiceMelhor])
		resposta.Individuo = populacao[indiceMelhor]

		fmt.Printf("Pontuacao: ")
		fmt.Println(scoreMelhor)
		resposta.Scores = append(resposta.Scores, scoreMelhor)

		// verifica se encontrou a solução ótima global
		//if(scoreMelhor == tamGenes)
		//	break;
	}

	c.JSON(http.StatusOK, resposta)
}

// Com base nos itens do orçamento, monta a população inicial com score para cada produto
func obterPopulacaoInicial(orcamento models.Orcamento, populacao *[]models.Individuo) {

	// Conexão com o banco de dados
	db := initDb()
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

	// Obtem o ultimo produto proposto para cada produto solicitado no orçamento
	var produtos []*models.PropostaProduto
	db.
		Where("cliente_produto_id in (?)", produtosOrcamento).
		Group("fornecedor_id,cliente_produto_id").
		Order("fornecedor_id desc, created_at desc").
		Find(&produtos)

	// Busca todos os fornecedores associados a esses orçamentos e que vão compor a população inicial
	var fornecedores []uint64
	db.
		Model(&models.OrcamentoFornecedor{}).
		Where("orcamento_id in (?)", orcamentos).
		Group("fornecedor_id").
		Pluck("fornecedor_id", &fornecedores)

	// Entregas Atrasadas
	var totais []Total
	var total Total
	rows, _ := db.
		Model(&models.Pedido{}).
		Select("fornecedor_id, sum(entrega_atrasada) atrasos, count(1) compras, timestampdiff(DAY, created_at, now()) tempo").
		Where("fornecedor_id in (?)", fornecedores).
		Group("fornecedor_id").
		Order("created_at").
		Rows()
	for rows.Next() {
		db.ScanRows(rows, &total)
		totais = append(totais, total)
	}

	// // Total de Pedidos
	// var atrasos []Total
	// rows, _ = db.
	// 	Model(&models.Pedido{}).
	// 	Select("fornecedor_id, count(1) as total").
	// 	Where("fornecedor_id in (?)", fornecedores).
	// 	Group("fornecedor_id").
	// 	Rows()
	// for rows.Next() {
	// 	var atraso Total
	// 	db.ScanRows(rows, &atraso)
	// 	atrasos = append(atrasos, atraso)
	// }

	// // Tempo de negociacao
	// db.
	// Model(&Pedido{}).
	// Select("fornecedor_id, count(1) as total").
	// Where("fornecedor_id in (?)", fornecedores).
	// Group("fornecedor_id")

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

		for _, t := range totais {
			if t.FornecedorID == fornecedorAtual {
				total = t
				break
			}
		}

		var individuo = models.Individuo{
			ID:      individuoID,
			Geracao: 1}

		// Adiciona novo individuo na populacao inicial
		*populacao = append(*populacao, individuo)

		// Adiciona os genes ao individuo
		for i := range orcamento.Produtos {

			found = false

			// Busca algum historico de orçamento para o produto e associa os dados da ultima proposta
			for j := range produtos {
				if orcamento.Produtos[i].ProdutoID == produtos[j].ClienteProdutoID && fornecedorAtual == produtos[j].FornecedorID {
					gene = models.Gene{
						OrcamentoProdutoID: i,
						ClienteProdutoID:   orcamento.Produtos[i].ProdutoID,
						IndividuoID:        individuoID,
						FornecedorID:       fornecedorAtual,
						Preco:              produtos[j].Preco,
						PrazoEntrega:       float32(produtos[j].PrazoEntrega),
						Atrasos:            total.Atrasos,
						Compras:            total.Compras}
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
			(*populacao)[individuoID].Genes = append((*populacao)[individuoID].Genes, gene)
		}
		individuoID++
	}
}

// Obtem o score da populacao inicial com base negociacoes do fornecedor
func gerarScorePopulacaoInicial(populacao *[]models.Individuo, solicitacao models.Solicitacao) {

	// Score por produto
	for i := range (*populacao)[0].Genes {

		// Preco
		var menorPreco, maiorPreco, notaPreco float32
		var menorPrazo, maiorPrazo, notaPrazo float32
		var menosAtrasos, maisAtrasos, notaAtrasos float32
		var menosCompras, maisCompras, notaCompras float32

		// Maior Preço
		for j, individuo := range *populacao {

			if j == 0 {
				menorPreco = individuo.Genes[i].Preco
				menorPrazo = individuo.Genes[i].PrazoEntrega
				menosAtrasos = individuo.Genes[i].Atrasos
				menosCompras = individuo.Genes[i].Compras
			}

			// Preço
			if individuo.Genes[i].Preco > maiorPreco {
				maiorPreco = individuo.Genes[i].Preco
			} else if individuo.Genes[i].Preco < menorPreco {
				menorPreco = individuo.Genes[i].Preco
			}

			// Prazo
			if individuo.Genes[i].PrazoEntrega > maiorPrazo {
				maiorPrazo = individuo.Genes[i].PrazoEntrega
			} else if individuo.Genes[i].PrazoEntrega < menorPrazo {
				menorPrazo = individuo.Genes[i].PrazoEntrega
			}

			// Atrasos
			if individuo.Genes[i].Atrasos > maisAtrasos {
				maisAtrasos = individuo.Genes[i].Atrasos
			} else if individuo.Genes[i].Atrasos < menosAtrasos {
				menosAtrasos = individuo.Genes[i].Atrasos
			}

			// Compras
			if individuo.Genes[i].Compras > maisCompras {
				maisCompras = individuo.Genes[i].Compras
			} else if individuo.Genes[i].Compras < menosCompras {
				menosCompras = individuo.Genes[i].Compras
			}
		}

		for j := range *populacao {
			notaPreco = (maiorPreco - (*populacao)[j].Genes[i].Preco) * (solicitacao.Preco / (maiorPreco - menorPreco))
			notaPrazo = (maiorPrazo - (*populacao)[j].Genes[i].PrazoEntrega) * (solicitacao.Prazo / (maiorPrazo - menorPrazo))
			notaAtrasos = (maisAtrasos - (*populacao)[j].Genes[i].Atrasos) * (solicitacao.Atrasadas / (maisAtrasos - menosAtrasos))
			notaCompras = (maisCompras - (*populacao)[j].Genes[i].Compras) * (solicitacao.Negociacoes / (maisCompras - menosCompras))
			(*populacao)[j].Genes[i].Score += uint64(notaPreco)
			(*populacao)[j].Genes[i].Score += uint64(notaAtrasos)
			(*populacao)[j].Genes[i].Score += uint64(notaCompras)
			(*populacao)[j].Genes[i].Score += uint64(notaPrazo)
		}
	}
}

// retorna o score do indivíduo
func obterPontuacao(individuo models.Individuo) (soma uint64) {

	// o score é a soma dos valores dos genes
	for i := range individuo.Genes {
		soma += individuo.Genes[i].Score
	}
	return soma
}

// realiza o cruzamento
func cruzamento(indicePai1 int, indicePai2 int, populacao []models.Individuo, filho *models.Individuo) {

	tamGenes := len(populacao[indicePai1].Genes)

	// escolhe um ponto aleatoriamente no intervalo [0, tamGenes - 1]
	ponto := rand.Intn(tamGenes)

	for i := 0; i < ponto; i++ {
		(*filho).Genes = append((*filho).Genes, populacao[indicePai1].Genes[i])
	}
	for i := ponto; i < tamGenes; i++ {
		(*filho).Genes = append((*filho).Genes, populacao[indicePai2].Genes[i])
	}
}

// realiza a mutação
func mutacao(individuo *models.Individuo, tamGenes int) {
	// escolhe um gene aleatoriamente no intervalo [0, tam_genes - 1]
	gene := rand.Intn(tamGenes)

	// modifica o valor do gene escolhido
	individuo.Genes[gene].Score++
}

// retorna o índice do melhor indivíduo da população
func obterMelhor(populacao []models.Individuo) (indiceMelhor int) {

	var score uint64
	scoreMelhor := obterPontuacao(populacao[0])

	for i := range populacao {
		score = obterPontuacao(populacao[i])
		if score > scoreMelhor {
			indiceMelhor = i
			scoreMelhor = score
		}
	}

	return indiceMelhor
}
