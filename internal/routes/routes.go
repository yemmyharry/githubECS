package routes

import (
	"github.com/gin-gonic/gin"
	"githubECS/internal/repository"
	"gorm.io/gorm"
	"net/http"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.POST("/discover", func(c *gin.Context) {
		var json struct {
			Query string `json:"query" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		repository.DiscoverRepos(json.Query, db)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})
}
