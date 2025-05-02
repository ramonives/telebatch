package extractor

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func ExtractByDate(date time.Time, db *sql.DB) ([]map[string]interface{}, error) {
	start := time.Now()

	query := `
		SELECT id, nome_cliente, telefone, data_ligacao,
		       duracao_segundos, motivo_ligacao, resolvido, atendente
		FROM chamadas_telemarketing
		WHERE DATE(data_ligacao) = $1
	`

	rows, err := db.Query(query, date.Format("2006-01-02"))
	if err != nil {
		return nil, fmt.Errorf("erro ao executar query: %w", err)
	}
	defer rows.Close()

	var results []map[string]interface{}

	for rows.Next() {
		var (
			id              int
			nomeCliente     string
			telefone        string
			dataLigacao     time.Time
			duracaoSegundos int
			motivoLigacao   string
			resolvido       bool
			atendente       string
		)

		err := rows.Scan(&id, &nomeCliente, &telefone, &dataLigacao,
			&duracaoSegundos, &motivoLigacao, &resolvido, &atendente)

		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha: %w", err)
		}

		results = append(results, map[string]interface{}{
			"id":               id,
			"nome_cliente":     nomeCliente,
			"telefone":         telefone,
			"data_ligacao":     dataLigacao,
			"duracao_segundos": duracaoSegundos,
			"motivo_ligacao":   motivoLigacao,
			"resolvido":        resolvido,
			"atendente":        atendente,
		})
	}

	duration := time.Since(start)

	// * Gera nome do lote
	lote := fmt.Sprintf("telemarketing_%s", time.Now().Format("02-01-2006_15:04:05"))

	// * Gera log
	logContent := fmt.Sprintf(
		"Lote: %s\nData da consulta: %s\nQuantidade: %d\nTempo: %s\n",
		lote, date.Format("2006-01-02"), len(results), duration.String(),
	)

	err = os.WriteFile(fmt.Sprintf("%s.log", lote), []byte(logContent), 0644)
	if err != nil {
		log.Printf("Erro ao gerar log: %v", err)
	} else {
		log.Printf("📝 Log criado: %s.log", lote)
	}

	return results, nil
}
