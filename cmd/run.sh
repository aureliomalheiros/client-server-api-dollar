#!/bin/bash

# Função para imprimir mensagem de log
log() {
    echo "$(date +'%Y-%m-%d %H:%M:%S') - $1"
}

log "Start..."

# Construir e iniciar os serviços Docker
docker-compose up --build -d

# Aguardar alguns segundos para garantir que o servidor esteja pronto
log "Aguardando o servidor iniciar..."
sleep 10

log "Run cliente Go..."
# Executar o cliente Go
docker-compose exec app go run client/client.go

# Opcional: Parar os serviços Docker após a execução do cliente
# log "Parando os containers Docker..."
# docker-compose down

log "Finishe..."
