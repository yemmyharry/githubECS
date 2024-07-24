package routes

import (
	"github.com/gin-gonic/gin"
	"githubECS/internal/handlers"
	"gorm.io/gorm"
)

func SetupRepoDiscoveryRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.POST("/search", handlers.SearchHandler)
	r.GET("/repositories/:full_name", handlers.GetRepositoryHandler)
	r.GET("/search", handlers.SearchByLanguageHandler)
	r.GET("/top", handlers.GetTopRepositoriesHandler)

	return r
}
