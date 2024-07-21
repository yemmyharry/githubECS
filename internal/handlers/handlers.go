package handlers

import (
	"githubECS/internal/repository"
	"githubECS/models"
	"gorm.io/gorm/clause"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRepositoryHandler(c *gin.Context, db *gorm.DB) {
	full_name := c.Param("full_name")
	full_name, err := url.QueryUnescape(full_name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository full_name"})
		return
	}

	var repos []models.Repository
	if err := db.Where("full_name LIKE ? OR description LIKE ?", "%"+full_name+"%", "%"+full_name+"%").Find(&repos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search repositories"})
		return
	}
	if len(repos) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func GetCommitsHandler(c *gin.Context, db *gorm.DB) {
	full_name := c.Param("full_name")
	full_name, err := url.QueryUnescape(full_name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repository full_name"})
		return
	}

	var repo models.Repository
	if err := db.Where("full_name LIKE ? OR description LIKE ?", "%"+full_name+"%", "%"+full_name+"%").First(&repo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Repository not found"})
		return
	}

	var commits []models.Commit
	if err := db.Where("repository_id = ?", repo.ID).Find(&commits).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve commits"})
		return
	}
	c.JSON(http.StatusOK, commits)
}

func SearchHandler(c *gin.Context, db *gorm.DB) {
	query := c.Query("query")
	repository.FindRepos(query, db)
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func SearchByLanguageHandler(c *gin.Context, db *gorm.DB) {
	language := c.Query("language")
	var repos []models.Repository
	if err := db.Where("LOWER(language) = ?", language).Find(&repos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search repositories"})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func GetTopRepositoriesHandler(c *gin.Context, db *gorm.DB) {
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
}

func ResetStartDateHandler(c *gin.Context, db *gorm.DB) {
	startDate := c.Query("start_date")
	if startDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Start date is required"})
		return
	}

	parsedDate, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
		return
	}

	if err := setStartDate(db, parsedDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update start date"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Start date updated", "start_date": startDate})
}

func setStartDate(db *gorm.DB, startDate time.Time) error {
	config := models.Config{
		Key:   "start_date",
		Value: startDate.Format(time.RFC3339),
	}
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&config).Error
}
