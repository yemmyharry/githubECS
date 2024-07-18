package models

import "time"

type Commit struct {
	SHA         string        `gorm:"primaryKey" json:"sha"`
	NodeID      string        `json:"node_id"`
	Commit      CommitDetails `gorm:"embedded;embeddedPrefix:commit_" json:"commit"`
	URL         string        `json:"url"`
	HTMLURL     string        `json:"html_url"`
	CommentsURL string        `json:"comments_url"`
	Author      *Person       `gorm:"embedded;embeddedPrefix:author_" json:"author"`
	Committer   *Person       `gorm:"embedded;embeddedPrefix:committer_" json:"committer"`
	Parents     []Parent      `gorm:"-" json:"parents"` // Exclude from GORM
}

type CommitDetails struct {
	Author       Person       `gorm:"embedded;embeddedPrefix:author_" json:"author"`
	Committer    Person       `gorm:"embedded;embeddedPrefix:committer_" json:"committer"`
	Message      string       `json:"message"`
	Tree         Tree         `gorm:"embedded;embeddedPrefix:tree_" json:"tree"`
	URL          string       `json:"url"`
	CommentCount int          `json:"comment_count"`
	Verification Verification `gorm:"embedded;embeddedPrefix:verification_" json:"verification"`
}

type Person struct {
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

type Tree struct {
	SHA string `json:"sha"`
	URL string `json:"url"`
}

type Verification struct {
	Verified  bool   `json:"verified"`
	Reason    string `json:"reason"`
	Signature string `json:"signature"`
	Payload   string `json:"payload"`
}

type Parent struct {
	SHA     string `json:"sha"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}
