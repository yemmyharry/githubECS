package routes

import (
	"github.com/gin-gonic/gin"
	"githubECS/internal/handlers"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/search", func(c *gin.Context) {
		handlers.SearchHandler(c, db)
	})

	r.GET("/repositories/:full_name", func(c *gin.Context) {
		handlers.GetRepositoryHandler(c, db)
	})

	r.GET("/repositories/:full_name/commits", func(c *gin.Context) {
		handlers.GetCommitsHandler(c, db)
	})

	r.GET("/search", func(c *gin.Context) {
		handlers.SearchByLanguageHandler(c, db)
	})

	r.GET("/top", func(c *gin.Context) {
		handlers.GetTopRepositoriesHandler(c, db)
	})

	r.POST("/reset_start_date", func(c *gin.Context) {
		handlers.ResetStartDateHandler(c, db)
	})

	return r
}
