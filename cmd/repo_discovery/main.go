package main

import (
	"githubECS/internal/routes"
	"githubECS/pkg/db"
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

	router := routes.SetupRepoDiscoveryRouter(database)
	log.Fatal(http.ListenAndServe(":8088", router))
}
