package commit

import (
	"encoding/json"
	"fmt"
	"githubECS/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
	"time"
)

func WatchCommits(db *gorm.DB) {
	var repos []models.Repository
	db.Find(&repos)

	for _, repo := range repos {
		log.Printf("Checking commits for repo: %s", repo.FullName)
		checkCommits(*repo.FullName, db, repo.ID)
	}
}

func checkCommits(name string, db *gorm.DB, repoID uint) {
	startTime, err := getStartDate(db)
	if err != nil {
		log.Printf("Error getting start date: %v", err)
		startTime = time.Time{}
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/commits?since=%s", name, startTime.Format(time.RFC3339))
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching commits:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	var commits []models.PartialCommit
	if err := json.Unmarshal(body, &commits); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	saveCommits(commits, db, repoID)
}

func saveCommits(commits []models.PartialCommit, db *gorm.DB, repoID uint) {
	for _, c := range commits {
		commit := models.Commit{
			Message:      c.Commit.Message,
			AuthorName:   c.Commit.Author.Name,
			AuthorEmail:  c.Commit.Author.Email,
			Date:         c.Commit.Author.Date,
			URL:          c.HTMLURL,
			RepositoryID: repoID,
		}

		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "url"}},
			DoNothing: true,
		}).Create(&commit)
		if result.Error != nil {
			log.Printf("Error saving commit %+v: %s", commit, result.Error)
		} else {
			log.Printf("Saved commit %+v", commit)
		}
	}
}

func getStartDate(db *gorm.DB) (time.Time, error) {
	var config models.Config
	if err := db.Where("key = ?", "start_date").First(&config).Error; err != nil {
		return time.Time{}, err
	}
	return time.Parse(time.RFC3339, config.Value)
}
