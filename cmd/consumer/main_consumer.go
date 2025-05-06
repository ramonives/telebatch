package main

import (
	"log"

	"github.com/ramonives/telebatch/internal/database"
	"github.com/ramonives/telebatch/internal/kafka"
)

func main() {
	// Conectar ao banco de dados
	database.Connect()

	// Garantir que a tabela de destino exista
	database.CreateImportedTable()

	log.Println("📥 Iniciando consumo de mensagens pendentes do Kafka...")

	// Iniciar o consumo
	kafka.ConsumeAndPersist(database.DB)
}
