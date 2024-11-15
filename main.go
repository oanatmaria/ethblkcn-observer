package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oanatmaria/ethblkcn-observer/client"
	"github.com/oanatmaria/ethblkcn-observer/parser"
	"github.com/oanatmaria/ethblkcn-observer/server"
	"github.com/oanatmaria/ethblkcn-observer/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		cancel()
	}()

	storage := storage.NewMemoryStorage()
	client := client.NewEthClient()
	parser, err := parser.NewEthParser(storage, client)
	if err != nil {
		log.Fatalf("Server error: can not start server, err: %v", err)
	}

	if parser == nil {
		log.Fatal("Server error: can not start server, failed to fetch latest block number")
	}

	server := server.NewHttpServer(":8080", parser)

	log.Println("Starting server...")
	if err := server.Start(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
