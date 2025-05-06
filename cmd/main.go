package main

import (
	"log"
	"time"

	"github.com/ramonives/telebatch/internal/database"
	"github.com/ramonives/telebatch/internal/extractor"
	"github.com/ramonives/telebatch/internal/kafka"
)

func main() {
	// * Conectar ao banco
	database.Connect()

	// * Criar tabela se necessário
	database.CreateTable()
	database.CreateImportedTable()

	// Definir data da extração (ex: ontem)
	dataExtracao := time.Now().AddDate(0, 0, -1)

	// Extrair registros
	registros, err := extractor.ExtractByDate(dataExtracao, database.DB)
	if err != nil {
		log.Fatalf("Erro na extração: %v", err)
	}

	log.Printf("Extração concluída com %d registros para a data %s",
		len(registros), dataExtracao.Format("2006-01-02"))

	if len(registros) == 0 {
		log.Println("Nenhum registro encontrado para enviar ao Kafka.")
		return
	}

	// Gerar nome do lote
	lote := "telemarketing_" + time.Now().Format("02-01-2006_15:04:05")

	// Converter para []any
	var registrosAny []any
	for _, r := range registros {
		registrosAny = append(registrosAny, r)
	}

	// Enviar para Kafka
	kafka.SendToKafka(lote, registrosAny)

	log.Println("Envio para o Kafka finalizado.")

	// ! Comentado para evitar conflito com cmd/consumer
	// kafka.ConsumeAndPersist(database.DB)
}
