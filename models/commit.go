package models

import "time"

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
