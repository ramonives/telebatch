package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

func SendToKafka(lote string, registros []any) {
	start := time.Now()
	ctx := context.Background()

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "telemarketing_lote",
	})
	defer writer.Close()

	// Enviar header
	headerMsg, _ := json.Marshal(map[string]interface{}{
		"type":      "header",
		"lote":      lote,
		"timestamp": time.Now().Format(time.RFC3339),
		"total":     len(registros),
	})
	writer.WriteMessages(ctx, kafka.Message{Value: headerMsg})

	// Enviar em batches de 1000
	batchSize := 1000
	for i := 0; i < len(registros); i += batchSize {
		end := i + batchSize
		if end > len(registros) {
			end = len(registros)
		}

		var batch []kafka.Message
		for _, r := range registros[i:end] {
			jsonMsg, _ := json.Marshal(r)
			batch = append(batch, kafka.Message{Value: jsonMsg})
		}

		if err := writer.WriteMessages(ctx, batch...); err != nil {
			log.Fatalf("❌ Erro ao enviar batch %d-%d para Kafka: %v", i+1, end, err)
		}

		log.Printf("🚚 Enviado batch %d - %d para o Kafka", i+1, end)
	}

	// Enviar footer
	footerMsg, _ := json.Marshal(map[string]interface{}{
		"type":      "footer",
		"lote":      lote,
		"timestamp": time.Now().Format(time.RFC3339),
		"total":     len(registros),
	})
	writer.WriteMessages(ctx, kafka.Message{Value: footerMsg})

	// Gerar log
	duration := time.Since(start)
	logContent := fmt.Sprintf(
		"Lote: %s\nTotal enviado: %d\nTempo: %s\nData/hora: %s\n",
		lote, len(registros), duration.String(), time.Now().Format("2006-01-02 15:04:05"),
	)

	err := os.WriteFile(fmt.Sprintf("%s-envio.log", lote), []byte(logContent), 0644)
	if err != nil {
		log.Printf("Erro ao gerar log de envio: %v", err)
	} else {
		log.Printf("📝 Log de envio gerado: %s-envio.log", lote)
	}
}
