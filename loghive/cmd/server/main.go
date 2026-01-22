package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexanderritik/loghive/internal/analyzer"
	"github.com/alexanderritik/loghive/internal/ingestion"
	"github.com/alexanderritik/loghive/internal/parser"
)

func main() {

	server := ingestion.NewServer(":3000", 1000)

	parser := parser.NewService(server.DataCh)

	analyser := analyzer.NewService(parser.ResultsCh)

	analyser.Start()

	parser.StartWorkers(5)

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("\nShutdown signal received...")

	server.Stop()

	log.Println("LogHive system exited successfully.")
}
