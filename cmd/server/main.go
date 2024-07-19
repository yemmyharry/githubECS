// package main
//
// import (
//
//	"github.com/gin-gonic/gin"
//	"github.com/joho/godotenv"
//	"githubECS/internal/routes"
//	"githubECS/internal/scheduler"
//	"githubECS/pkg/db"
//	"log"
//	"os"
//
// )
//
//	func main() {
//		if err := godotenv.Load(); err != nil {
//			log.Fatal("Error loading .env file")
//		}
//
//		dsn := os.Getenv("DATABASE_URL")
//		database, err := db.Initialize(dsn)
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		scheduler.StartScheduler(database)
//
//		r := gin.Default()
//		routes.SetupRouter(database)
//
//		err = r.Run(":8089")
//		if err != nil {
//			return
//		}
//	}
package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"githubECS/internal/routes"
	"githubECS/pkg/db"
)

func main() {
	// Load environment variables from .env file if it exists
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	database, err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}

	router := routes.SetupRouter(database)
	router.Run(":8080")
}
