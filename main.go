package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aureliomalheiros/client-server-api-dollar/server"
)

func main() {
	// Iniciando o servidor
	go func() {
		err := server.StartServer()
		if err != nil {
			log.Fatalf("Erro ao iniciar o servidor: %v", err)
		}
	}()

	// Aguardando um sinal para encerrar
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Encerrando o servidor
	err := server.StopServer()
	if err != nil {
		log.Fatalf("Erro ao encerrar o servidor: %v", err)
	}

	fmt.Println("Servidor encerrado.")
}
