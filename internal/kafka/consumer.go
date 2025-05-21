package kafka

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

type Ligacao struct {
	ID              int       `json:"id"`
	NomeCliente     string    `json:"nome_cliente"`
	Telefone        string    `json:"telefone"`
	DataLigacao     time.Time `json:"data_ligacao"`
	DuracaoSegundos int       `json:"duracao_segundos"`
	MotivoLigacao   string    `json:"motivo_ligacao"`
	Resolvido       bool      `json:"resolvido"`
	Atendente       string    `json:"atendente"`
	Type            string    `json:"type"`
	Lote            string    `json:"lote"`
}

func ConsumeAndPersist(db *sql.DB) {
	start := time.Now()
	ctx := context.Background()

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "telemarketing_lote",
		GroupID:  "telebatch_consumer_batch10k",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("📥 Iniciando consumo do Kafka...")

	count := 0
	var lote string
	var buffer []Ligacao
	batchSize := 50000

	insertBatch := func(batch []Ligacao) {
		if len(batch) == 0 {
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Erro ao iniciar transação: %v", err)
			return
		}

		stmt, err := tx.Prepare(`
			INSERT INTO chamadas_importadas
			(nome_cliente, telefone, data_ligacao, duracao_segundos, motivo_ligacao, resolvido, atendente)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`)
		if err != nil {
			log.Printf("Erro ao preparar statement: %v", err)
			return
		}
		defer stmt.Close()

		for _, msg := range batch {
			_, err := stmt.Exec(msg.NomeCliente, msg.Telefone, msg.DataLigacao, msg.DuracaoSegundos, msg.MotivoLigacao, msg.Resolvido, msg.Atendente)
			if err != nil {
				log.Printf("Erro ao inserir no batch: %v", err)
			}
		}

		if err := tx.Commit(); err != nil {
			log.Printf("Erro ao commitar batch: %v", err)
		}
	}

	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}

		var msg Ligacao
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Erro ao fazer unmarshal: %v", err)
			continue
		}

		if msg.Type == "header" {
			if msg.Lote != "" {
				lote = msg.Lote
			}
			continue
		}

		if msg.Type == "footer" {
			// Flush final do que sobrou no buffer
			insertBatch(buffer)
			break
		}

		if msg.Lote != "" && lote == "" {
			lote = msg.Lote
		}

		buffer = append(buffer, msg)
		count++

		if len(buffer) >= batchSize {
			insertBatch(buffer)
			log.Printf("✅ %d registros consumidos e inseridos...", count)
			buffer = buffer[:0] // limpa o buffer
		}
	}

	duration := time.Since(start)
	logContent := fmt.Sprintf(
		"Lote: %s\nTotal processado: %d\nTempo: %s\nData/hora: %s\n",
		lote, count, duration.String(), time.Now().Format("2006-01-02 15:04:05"),
	)

	filename := fmt.Sprintf("%s-consumo.log", lote)
	if lote == "" {
		filename = fmt.Sprintf("lote_undefined_%s-consumo.log", time.Now().Format("15_04_05"))
	}

	err := os.WriteFile(filename, []byte(logContent), 0644)
	if err != nil {
		log.Printf("Erro ao criar log de consumo: %v", err)
	} else {
		log.Printf("📝 Log de consumo gerado: %s", filename)
	}
}
