package handlers

import (
	"encoding/json"
	"githubECS/internal/commit"
	"githubECS/internal/repository"
	"githubECS/models"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func getDB(c *gin.Context) *gorm.DB {
	return c.MustGet("db").(*gorm.DB)
}

func SearchHandler(c *gin.Context) {
	db := getDB(c)
	query := c.Query("query")
	repository.FindRepos(query, db)
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func GetRepositoryHandler(c *gin.Context) {
	db := getDB(c)
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

func GetCommitsHandler(c *gin.Context) {
	db := getDB(c)
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

func SearchByLanguageHandler(c *gin.Context) {
	db := getDB(c)
	language := c.Query("language")
	var repos []models.Repository
	if err := db.Where("LOWER(language) = ?", language).Find(&repos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search repositories"})
		return
	}
	c.JSON(http.StatusOK, repos)
}

func GetTopRepositoriesHandler(c *gin.Context) {
	db := getDB(c)
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

func ResetStartDateHandler(c *gin.Context) {
	db := getDB(c)
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

	if err := models.SetStartDate(db, parsedDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update start date"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Start date updated", "start_date": startDate})
}

func StartCommitManagerConsumer(db *gorm.DB, rabbitCh *amqp.Channel) {
	queueName := "commit_manager_queue"
	msgs, err := rabbitCh.Consume(
		queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var req struct {
				Repo string `json:"repo"`
			}
			if err := json.Unmarshal(d.Body, &req); err != nil {
				log.Printf("Error unmarshalling JSON: %v", err)
				continue
			}
			commit.WatchCommits(db)
		}
	}()

	log.Printf("Waiting for messages. To exit press CTRL+C")
	<-forever
}
