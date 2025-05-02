package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	connStr := "host=localhost port=5432 user=admin password=admin dbname=telebatch sslmode=disable"

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o banco: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erro ao testar conexão com o banco: %v", err)
	}

	fmt.Println("Conexão com PostgreSQL foi um sucesso!")
}

func CreateTable() {
	query := `
	CREATE TABLE IF NOT EXISTS chamadas_telemarketing (
		id SERIAL PRIMARY KEY,
		nome_cliente VARCHAR(100),
		telefone VARCHAR(20),
		data_ligacao TIMESTAMP,
		duracao_segundos INT,
		motivo_ligacao VARCHAR(100),
		resolvido BOOLEAN,
		atendente VARCHAR(100)
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	fmt.Println("Tabela 'chamadas_telemarketing' verificada/criada.")
}
