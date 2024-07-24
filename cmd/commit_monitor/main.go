package main

import (
	"githubECS/internal/scheduler"
	"githubECS/pkg/db"
	"githubECS/rabbitmq"
	"log"
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

	scheduler.StartCommitMonitor(database, rabbitCh)
}
