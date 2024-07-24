package main

import (
	"githubECS/internal/handlers"
	"githubECS/internal/routes"
	"githubECS/pkg/db"
	"githubECS/rabbitmq"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	dsn := os.Getenv("DATABASE_URL")

	database, err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}

	rabbitConn, err := rabbitmq.NewConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitConn.Close()

	rabbitCh, err := rabbitConn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbitCh.Close()

	go handlers.StartCommitManagerConsumer(database, rabbitCh)

	router := routes.SetupCommitManagerRouter(database)
	log.Fatal(http.ListenAndServe(":8087", router))
}
