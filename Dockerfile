# Dockerfile para o projeto client-server-api-dollar

# Imagem base
FROM golang:latest

# Definir o diretório de trabalho
WORKDIR /go/src/client-server-api-dollar

# Copiar os arquivos go.mod e go.sum
COPY go.mod go.sum ./

# Baixar as dependências do Go
RUN go mod download

# Copiar o restante dos arquivos do projeto para o container
COPY . .

# Compilar o código
RUN go build -o main .

# Comando padrão para executar o servidor
CMD ["./main"]
