package main

import (
	"encoding/json"
	"githubECS/models"
	"io"
	"log"

	"gorm.io/gorm/clause"
	"net/http"
)

const commitURL = "https://api.github.com/repos/"

func watchCommits() {
	var repos []models.Repository
	db.Find(&repos)

	for _, repo := range repos {
		log.Printf("Checking commits for repo: %s", repo.FullName)
		checkCommits(repo.FullName)
	}
}

func checkCommits(fullName string) {
	resp, err := http.Get(commitURL + fullName + "/commits")
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

	log.Printf("API response for repo %s: %s", fullName, body)

	var commits []models.Commit
	if err := json.Unmarshal(body, &commits); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		log.Println("Response body:", string(body))
		return
	}

	log.Printf("Parsed commits for repo %s: %+v", fullName, commits)

	saveCommits(commits)
}

func saveCommits(commits []models.Commit) {
	for _, commit := range commits {
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "sha"}},
			DoNothing: true,
		}).Create(&commit)
		if result.Error != nil {
			log.Printf("Error saving commit %+v: %s", commit, result.Error)
		} else {
			log.Printf("Saved commit %+v", commit)
		}
	}
}
