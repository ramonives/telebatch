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
		GroupID:  "telebatch_consumer_group_final",
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer reader.Close()

	log.Println("📥 Iniciando consumo do Kafka...")

	count := 0
	var lote string

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
			continue
		}

		if msg.Type == "footer" {
			break
		}

		// ✅ Só define o lote se estiver preenchido
		if msg.Lote != "" {
			lote = msg.Lote
		}

		_, err = db.Exec(`
			INSERT INTO chamadas_importadas
			(nome_cliente, telefone, data_ligacao, duracao_segundos, motivo_ligacao, resolvido, atendente)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`, msg.NomeCliente, msg.Telefone, msg.DataLigacao, msg.DuracaoSegundos, msg.MotivoLigacao, msg.Resolvido, msg.Atendente)

		if err != nil {
			log.Printf("Erro ao inserir no banco: %v", err)
			continue
		}

		count++
		if count%100 == 0 {
			log.Printf("✅ %d registros consumidos e inseridos...", count)
		}
	}

	duration := time.Since(start)

	logContent := fmt.Sprintf(
		"Lote: %s\nTotal processado: %d\nTempo: %s\nData/hora: %s\n",
		lote, count, duration.String(), time.Now().Format("2006-01-02 15:04:05"),
	)

	// ✅ Garante nome válido de arquivo
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
