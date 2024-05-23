FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o myapp main.go

CMD ["./client-server-api-dollar"]
