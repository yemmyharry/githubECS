package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"githubECS/internal/routes"
	"githubECS/internal/scheduler"
	"githubECS/pkg/db"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	database, err := db.Initialize(dsn)
	if err != nil {
		log.Fatalln(err)
	}

	scheduler.StartScheduler(database)

	r := gin.Default()
	routes.SetupRouter(r, database)

	err = r.Run()
	if err != nil {
		return
	}
}
