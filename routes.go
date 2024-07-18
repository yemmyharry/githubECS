package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.POST("/discover", func(c *gin.Context) {
		var json struct {
			Query string `json:"query" binding:"required"`
		}
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		DiscoverRepos(json.Query)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	return r
}
