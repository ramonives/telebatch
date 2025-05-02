# ✅ Desafio Monitoria – Pipeline de Telemarketing com Golang + Kafka + PostgreSQL

## 🔹 Fase 1 – Preparação do ambiente

- [x] Criar `docker-compose.yml` com:
  - [x] PostgreSQL (porta 5432)
  - [x] Kafka + Zookeeper (porta 9092)
- [x] Subir os containers e testar conexão com banco

- [x] Criar tabela `chamadas_telemarketing` com estrutura:
  - `id`, `nome_cliente`, `telefone`, `data_ligacao`, `duracao_segundos`, `motivo_ligacao`, `resolvido`, `atendente`

- [x] Gerar script para popular 5 milhões de registros variados
  - Datas variadas (últimos 30 dias)
  - Motivos e atendentes diversos

---

## 🔹 Fase 2 – Extração dos dados

- [x] Iniciar projeto Go com `go mod init`
- [x] Conectar ao banco PostgreSQL
- [x] Implementar extração de registros por `data_ligacao`
- [x] Medir:
  - [x] Tempo de extração
  - [x] Quantidade de registros
- [x] Nomear lote: `telemarketing_dd-mm-yyyy_hh-mm-ss`
- [x] Gerar log `.log` com essas informações

---

## 🔹 Fase 3 – Envio para o Kafka

- [ ] Criar produtor Kafka em Go
- [ ] Criar/usar tópico `telemarketing_lote`
- [ ] Enviar:
  - [ ] Header com início e total
  - [ ] Mensagens dos registros
  - [ ] Footer com fim e total entregue
- [ ] Medir tempo de envio
- [ ] Escrever log com tempo e totais enviados

---

## 🔹 Fase 4 – Consumo e persistência

- [ ] Criar tabela `chamadas_importadas` (mesma estrutura)
- [ ] Criar consumidor Kafka em Go
- [ ] Ignorar header/footer no consumo
- [ ] Inserir dados recebidos na nova tabela
- [ ] Validar que a quantidade bate com a extraída

---

## 🔹 Fase 5 – Validação Final

- [ ] Conferir logs:
  - Extração
  - Envio
  - Recebimento
- [ ] Conferir contagem nas tabelas:
  - `chamadas_telemarketing` (origem)
  - `chamadas_importadas` (destino)
