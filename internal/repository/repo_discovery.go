package repository

import (
	"encoding/json"
	"githubECS/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	searchURL         = "https://api.github.com/search/repositories?q="
	rateLimitResetURL = "https://api.github.com/rate_limit"
)

type RepoSearchResult struct {
	Items []models.Repository `json:"items"`
}

func DiscoverRepos(query string, db *gorm.DB) {
	if !checkRateLimit() {
		return
	}

	resp, err := http.Get(searchURL + query)
	if err != nil {
		log.Println("Error fetching repositories:", err)
		return
	}
	defer resp.Body.Close()

	handleRateLimitHeaders(resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
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

	saveRepos(itemsToSave, db)
}

func saveRepos(repos []models.Repository, db *gorm.DB) {
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

func checkRateLimit() bool {
	resp, err := http.Get(rateLimitResetURL)
	if err != nil {
		log.Println("Error checking rate limit:", err)
		return false
	}
	defer resp.Body.Close()

	var rateLimit struct {
		Resources struct {
			Core struct {
				Remaining int `json:"remaining"`
				Reset     int `json:"reset"`
			} `json:"core"`
		} `json:"resources"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rateLimit); err != nil {
		log.Println("Error decoding rate limit response:", err)
		return false
	}

	if rateLimit.Resources.Core.Remaining == 0 {
		waitUntil := time.Unix(int64(rateLimit.Resources.Core.Reset), 0)
		waitDuration := time.Until(waitUntil)
		log.Printf("Rate limit exceeded. Waiting for %v until %v", waitDuration, waitUntil)
		time.Sleep(waitDuration)
		return false
	}

	return true
}

func handleRateLimitHeaders(headers http.Header) {
	remaining, err := strconv.Atoi(headers.Get("X-RateLimit-Remaining"))
	if err != nil {
		log.Println("Error parsing X-RateLimit-Remaining header:", err)
		return
	}

	if remaining == 0 {
		reset, err := strconv.Atoi(headers.Get("X-RateLimit-Reset"))
		if err != nil {
			log.Println("Error parsing X-RateLimit-Reset header:", err)
			return
		}

		resetTime := time.Unix(int64(reset), 0)
		waitDuration := time.Until(resetTime)
		log.Printf("Rate limit exceeded. Waiting for %v until %v", waitDuration, resetTime)
		time.Sleep(waitDuration)
	}
}
