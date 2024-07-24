package routes

import (
	"github.com/gin-gonic/gin"
	"githubECS/internal/handlers"
	"gorm.io/gorm"
)

func SetupCommitManagerRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/repositories/:full_name/commits", handlers.GetCommitsHandler)
	r.POST("/reset_start_date", handlers.ResetStartDateHandler)

	return r
}
