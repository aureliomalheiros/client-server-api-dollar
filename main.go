package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/aureliomalheiros/client-server-api-dollar/client"
	"github.com/aureliomalheiros/client-server-api-dollar/server"
)

func main() {
	// Iniciar o servidor em uma goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	client.Run()

	// Esperar por sinais de interrupção para desligar o servidor de forma graciosa
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server...")
}
