package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/online-tryout/parsing-sheets-api/api"
	"github.com/online-tryout/parsing-sheets-api/broker"
	"github.com/online-tryout/parsing-sheets-api/util"
)

// @title Parsing Sheet API Documentation
// @version 1.0
// @description This is a documentation for Online Tryout Apps

// @host localhost:8081
// @BasePath /
func main() {
	// configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config: ", err)
	}

	// rabbitmq
	rabbitmq, err := broker.NewRabbitMq(config.RabbitSource, &config)
	if err != nil {
		log.Fatal("can't connect to rabbitmq: ", err)
	}

	// Create a channel to signal when the server is ready to shutdown
	shutdown := make(chan struct{})

	// Start consuming messages in a separate goroutine
	go func() {
		err := rabbitmq.ConsumeEvent("parsing-sheets-queue")
		if err != nil {
			log.Fatalf("Failed to consume messages: %v", err)
		}
	}()

	// server
	server, err := api.NewServer(&config, rabbitmq)
	if err != nil {
		log.Fatal("can't create server: ", err)
	}

	// Start server
	go func() {
		err := server.Start(config.ServerAddress)
		if err != nil {
			log.Fatal("can't start server: ", err)
		}
	}()

	// Wait for a signal to shutdown
	<-shutdown
}
