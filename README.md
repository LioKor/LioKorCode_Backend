# LioKorCode Backend (REST API + WS)

## How to build:
* for GoLang 1.18.1+
* install Redis
* install PostgreSQL
* create database and user:
  * `psql -U postgres`
  * `CREATE DATABASE liokoredu;`
  * `CREATE USER lk WITH PASSWORD 'password';`
  * `GRANT ALL PRIVILEGES ON DATABASE liokoredu TO lk;`
* import `db.sql` or `schema.sql`:
  * `psql -U lk -d liokoredu -f schema.sql`
* install easyjson and create models
  * `go install github.com/mailru/easyjson/...@latest`
  * `cd application/models && easyjson -all ./`
* `go build cmd/main.go`

Backend for LioKorCode project made for VK Education | Technopark in BMSTU. 
Spring 2022.