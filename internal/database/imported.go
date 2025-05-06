package database

import (
	"log"
)

// CreateImportedTable garante que a tabela 'chamadas_importadas' exista
func CreateImportedTable() {
	query := `
	CREATE TABLE IF NOT EXISTS chamadas_importadas (
		id SERIAL PRIMARY KEY,
		nome_cliente TEXT,
		telefone TEXT,
		data_ligacao TIMESTAMP,
		duracao_segundos INTEGER,
		motivo_ligacao TEXT,
		resolvido BOOLEAN,
		atendente TEXT
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela chamadas_importadas: %v", err)
	}

	log.Println("Tabela 'chamadas_importadas' verificada/criada.")
}
