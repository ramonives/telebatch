package main

import (
	"log"
	"time"

	"github.com/ramonives/telebatch/internal/database"
	"github.com/ramonives/telebatch/internal/extractor"
)

func main() {
	database.Connect()
	database.CreateTable()

	// ! database.PopulateTable(10000, database.DB)

	// ! log.Println("Inserção de teste concluída.")

	// * Realizar extração para uma data (ex: ontem)
	dataExtracao := time.Now().AddDate(0, 0, -1)

	resultado, err := extractor.ExtractByDate(dataExtracao, database.DB)
	if err != nil {
		log.Fatalf("Erro na extração: %v", err)
	}

	log.Printf("✅ Extração concluída com %d registros para a data %s", len(resultado), dataExtracao.Format("2006-01-02"))
}
