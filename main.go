package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/online-tryout/parsing-sheets-api/api"
	"github.com/online-tryout/parsing-sheets-api/broker"
	db "github.com/online-tryout/parsing-sheets-api/db/sqlc"
	"github.com/online-tryout/parsing-sheets-api/util"
)

func main() {
    // configuration
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config: ", err)
	}

	// postgresql
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't connect to database: ", err)
	}
	store := db.NewStore(conn)

	// rabbitmq
	rabbitmq, err := broker.NewRabbitMq(config.RabbitSource)
	if err != nil {
		log.Fatal("can't connect to rabbitmq: ", err)
	}

	// server
	server, err := api.NewServer(&config, store, rabbitmq)
	if err != nil {
		log.Fatal("can't create server: ", err)
	}

	// start server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server: ", err)
	}
}