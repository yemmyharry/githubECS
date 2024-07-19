package commit

import (
	"encoding/json"
	"githubECS/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
)

const commitURL = "https://api.github.com/repos/"

func WatchCommits(db *gorm.DB) {
	var repos []models.Repository
	db.Find(&repos)

	for _, repo := range repos {
		log.Printf("Checking commits for repo: %s", repo.FullName)
		checkCommits(repo.FullName, db, repo.ID)
	}
}

func checkCommits(name string, db *gorm.DB, repoID uint) {
	resp, err := http.Get(commitURL + name + "/commits")
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
