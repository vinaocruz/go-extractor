# GO Extractor

Este código apresenta um desafio técnico para extrair dados de arquivos zip e retorná-lo através de uma api de consulta

## Versões utilizadas

Para este projeto, foram utilizadas as seguintes versões:

* Docker 26.1
* Docker-compose 1.29
* Go 1.22

## Setup

Execute o docker-compose para subir o banco de dados postgres utilizado neste projeto. Ao criar o container, o sql da pasta initdb/ será executado no setup

```
docker-compose up -d
```

Realize uma cópia do arquivo .env.example para .env com dados de configuração da base

## Importação de dados

A base de dados será populada com os arquivos estáticos localizado na pasta `storage/example/` contendo +63M de registros. Para executar a importação do histórico de negociações da B3 dos últimos 7 dias, utilize:

```
go run . dataImport
```

**TODO**
* ler todos arquivos da pasta de exemplo (remover hard code)
* implementar um crawler para baixar os últimos dias direto do site da B3

### Estratégia de otimização

Para otimizar a importação da massa de dados, foram utilizadas as seguintes estratégias:

* Leitura do arquivo por linha, evitando carregá-lo todo em memória
* Goroutine para ler os arquivos em concorrência
* Channel para enviar lotes de dados lidos e prontos para inserção no banco
* Tabela inicializada sem índices
* Inserção em lote no banco com COPY FROM
* Ao concluir importação, criação de indices para otimizar consultas da API

### Estratégia de modelagem do banco

Para este código foi mantido apenas uma tabela, considerando que o código do ticket é uma chave que represente a empresa que opera na bolsa.

Foram criados dois índices para otimizar busca em `ticketCode` e `transationAt`

## Servidor Web

A API permitirá realizar dados na massa importada. Para inicia-lo, execute:

```
go run .
```

O servidor web é iniciado na porta 8080, com Gin framework. Para realizar consulta, acesse (exemplos):

http://localhost:8080/negociations?ticker=INDQ24

O segundo parâmetro opcional para recorte de data:

http://localhost:8080/negociations?ticker=PETR4&DataNegocio=2024-07-01