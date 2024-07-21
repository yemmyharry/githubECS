package main

import (
	"githubECS/internal/scheduler"
	"log"
	"os"

	"github.com/joho/godotenv"
	"githubECS/internal/routes"
	"githubECS/pkg/db"
)

func main() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	database, err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}
	scheduler.StartScheduler(database)

	router := routes.SetupRouter(database)
	router.Run(":8080")
}
