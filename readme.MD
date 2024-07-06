# GO Extractor

Este código apresenta um desafio técnico para extrair dados de arquivos zip e retorna-lo através de uma api de consulta

## Docker

Execute o docker-compose da raiz, para subir o banco de dados utilizado neste projeto

```
docker-compose up -d
```

## Importação de dados

Execute o command para baixar e importar o histórico de negociações da B3 dos últimos 7 dias

```
go run . dataImport
```

## Servidor Web

Para subir o servidor web, execute:

```
go run .
```