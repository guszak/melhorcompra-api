# melhorcompra-api

API - Aplicação para predição de compras baseado em Algoritmo Genético


## Configurar variáveis de ambiente (opcional)

```
export MELHORCOMPRA_DATABASEHOST="localhost"
export MELHORCOMPRA_DATABASEPORT="3306"
export MELHORCOMPRA_DATABASEUSER="root"
export MELHORCOMPRA_DATABASEPASS="root"
export MELHORCOMPRA_DATABASENAME="melhorcompra"
export MELHORCOMPRA_PORT="3001"
```

## Instalar dependências

```bash
go get -t ./...
```

## Instalar dependências usando o Dep

Obtenha a gerenciador de pacotes via
```bash
go get -u github.com/golang/dep/cmd/dep
```
Instalando as dependências
```bash
dep ensure
```

## Build

```bash
go build
```