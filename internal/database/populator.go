package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
)

var motivos = []string{
	"Suporte técnico",
	"Cobrança",
	"Troca de plano",
	"Cancelamento",
	"Informações gerais",
}

var atendentes = []string{
	"Ana Paula", "Carlos Souza", "Fernanda Lima",
	"João Oliveira", "Marcos Dias", "Luciana Silva",
	"Beatriz Santos", "Rodrigo Melo", "Camila Torres", "Paulo Martins",
}

type Chamada struct {
	NomeCliente     string
	Telefone        string
	DataLigacao     time.Time
	DuracaoSegundos int
	MotivoLigacao   string
	Resolvido       bool
	Atendente       string
}

func GenerateFakeCall() Chamada {
	now := time.Now()
	diaAleatorio := rand.Intn(30)
	dataLigacao := now.AddDate(0, 0, -diaAleatorio)

	return Chamada{
		NomeCliente:     fmt.Sprintf("Cliente %d", rand.Intn(100000)),
		Telefone:        fmt.Sprintf("1%09d", rand.Intn(999999999)),
		DataLigacao:     dataLigacao,
		DuracaoSegundos: rand.Intn(1771) + 30,
		MotivoLigacao:   motivos[rand.Intn(len(motivos))],
		Resolvido:       rand.Intn(2) == 1,
		Atendente:       atendentes[rand.Intn(len(atendentes))],
	}
}

func PopulateTable(qtd int, db *sql.DB) {
	const batchSize = 1000
	insertBase := `INSERT INTO chamadas_telemarketing
	(nome_cliente, telefone, data_ligacao, duracao_segundos, motivo_ligacao, resolvido, atendente) VALUES `

	valCount := 0
	valueStrings := []string{}
	valueArgs := []interface{}{}

	start := time.Now()

	for i := 0; i < qtd; i++ {
		c := GenerateFakeCall()

		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			valCount+1, valCount+2, valCount+3, valCount+4, valCount+5, valCount+6, valCount+7))

		valueArgs = append(valueArgs,
			c.NomeCliente, c.Telefone, c.DataLigacao, c.DuracaoSegundos,
			c.MotivoLigacao, c.Resolvido, c.Atendente)

		valCount += 7

		if (i+1)%batchSize == 0 || i+1 == qtd {
			query := insertBase + strings.Join(valueStrings, ",")
			_, err := db.Exec(query, valueArgs...)
			if err != nil {
				log.Fatalf("Erro ao inserir lote no banco: %v", err)
			}

			log.Printf("Inseridos %d registros...", i+1)

			valueStrings = []string{}
			valueArgs = []interface{}{}
			valCount = 0
		}
	}

	elapsed := time.Since(start)
	log.Printf("Inserção de %d registros concluída em %s", qtd, elapsed)
}
