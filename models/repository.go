package models

import "time"

type Repository struct {
	ID              uint      `gorm:"primaryKey"`
	FullName        string    `gorm:"uniqueIndex" json:"full_name"`
	Description     string    `json:"description"`
	URL             string    `json:"html_url"`
	Language        string    `json:"language"`
	ForksCount      int       `json:"forks_count"`
	StarsCount      int       `json:"stargazers_count"`
	OpenIssuesCount int       `json:"open_issues_count"`
	WatchersCount   int       `json:"watchers_count"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
