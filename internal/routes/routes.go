package routes

import (
	"githubECS/internal/repository"
	"githubECS/models"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.POST("/search", func(c *gin.Context) {
		query := c.Query("query")
		repository.DiscoverRepos(query, db)
		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	})

	r.GET("/repositories/:full_name", func(c *gin.Context) {
		full_name := c.Param("full_name")
		full_name, err := url.QueryUnescape(full_name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository full_name"})
			return
		}

		var repos []models.Repository
		if err := db.Where("full_name ILIKE ? OR description ILIKE ?", "%"+full_name+"%", "%"+full_name+"%").Find(&repos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search repositories"})
			return
		}
		if len(repos) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
			return
		}
		c.JSON(http.StatusOK, repos)
	})

	r.GET("/repositories/:full_name/commits", func(c *gin.Context) {
		full_name := c.Param("full_name")
		full_name, err := url.QueryUnescape(full_name)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository full_name"})
			return
		}

		var repo models.Repository
		if err := db.Where("full_name ILIKE ? OR description ILIKE ?", "%"+full_name+"%", "%"+full_name+"%").First(&repo).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
			return
		}

		var commits []models.Commit
		if err := db.Where("repository_id = ?", repo.ID).Find(&commits).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve commits"})
			return
		}
		c.JSON(http.StatusOK, commits)
	})

	r.GET("/search", func(c *gin.Context) {
		language := c.Query("language")
		var repos []models.Repository
		if err := db.Where("LOWER(language) = ?", language).Find(&repos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search repositories"})
			return
		}
		c.JSON(http.StatusOK, repos)
	})

	r.GET("/top", func(c *gin.Context) {
		nStr := c.Query("n")
		n, err := strconv.Atoi(nStr)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number of repositories"})
			return
		}
		var repos []models.Repository
		if err := db.Order("stars_count desc").Limit(n).Find(&repos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top repositories"})
			return
		}
		c.JSON(http.StatusOK, repos)
	})
}
