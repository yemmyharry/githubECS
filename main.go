package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"githubECS/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	if err := db.AutoMigrate(&models.Repository{}, &models.Commit{}); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(1).Minutes().Do(watchCommits)

	scheduler.StartAsync()

	r := SetupRouter(db)

	err = r.Run(":8088")
	if err != nil {
		fmt.Printf("Error running %s:", err)
		return
	}
}
