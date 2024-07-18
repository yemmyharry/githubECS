package main

import (
	"encoding/json"
	"githubECS/models"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
)

const searchURL = "https://api.github.com/search/repositories?q="

type RepoSearchResult struct {
	Items []models.Repository `json:"items"`
}

func DiscoverRepos(query string) {
	resp, err := http.Get(searchURL + query)
	if err != nil {
		log.Println("Error fetching repositories:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	var result RepoSearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		log.Println("Error unmarshalling JSON:", err)
		return
	}

	itemsToSave := result.Items
	if len(itemsToSave) > 10 {
		itemsToSave = itemsToSave[:10]
	}

	saveRepos(itemsToSave)
}

func saveRepos(repos []models.Repository) {
	for _, repo := range repos {
		result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "full_name"}},
			DoNothing: true,
		}).Create(&repo)
		if result.Error != nil {
			log.Printf("Error saving repository %+v: %s", repo, result.Error)
		} else {
			log.Printf("Saved repository %+v", repo)
		}
	}
}
