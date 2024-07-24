package models

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"time"
)

type Repository struct {
	ID              uint       `gorm:"primaryKey"`
	FullName        *string    `gorm:"uniqueIndex" json:"full_name"`
	Description     *string    `json:"description"`
	URL             *string    `json:"html_url"`
	Language        *string    `json:"language"`
	ForksCount      *int       `json:"forks_count"`
	StarsCount      *int       `json:"stargazers_count"`
	OpenIssuesCount *int       `json:"open_issues_count"`
	WatchersCount   *int       `json:"watchers_count"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}

type Commit struct {
	ID           uint      `gorm:"primaryKey"`
	Message      string    `json:"message"`
	AuthorName   string    `json:"author_name"`
	AuthorEmail  string    `json:"author_email"`
	Date         time.Time `json:"date"`
	URL          string    `gorm:"uniqueIndex" json:"url"`
	RepositoryID uint      `json:"repository_id"`
}

type PartialCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
		Author  struct {
			Name  string    `json:"name"`
			Email string    `json:"email"`
			Date  time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	HTMLURL string `json:"html_url"`
}

type Config struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SetStartDate(db *gorm.DB, startDate time.Time) error {
	config := Config{
		Key:   "start_date",
		Value: startDate.Format(time.RFC3339),
	}
	return db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&config).Error
}
