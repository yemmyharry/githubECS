package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	switch v := value.(type) {
	case string:
		*a = strings.Split(v, ",")
	case []byte:
		*a = strings.Split(string(v), ",")
	default:
		return fmt.Errorf("unsupported data type: %T", v)
	}
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "", nil
	}
	return strings.Join(a, ","), nil
}

type Repository struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	Name            string      `json:"name"`
	FullName        string      `gorm:"unique" json:"full_name"`
	Private         bool        `json:"private"`
	OwnerLogin      string      `json:"owner.login"`
	OwnerID         int         `json:"owner.id"`
	OwnerAvatarURL  string      `json:"owner.avatar_url"`
	HTMLURL         string      `json:"html_url"`
	Description     string      `json:"description"`
	Fork            bool        `json:"fork"`
	URL             string      `json:"url"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	PushedAt        time.Time   `json:"pushed_at"`
	GitURL          string      `json:"git_url"`
	SSHURL          string      `json:"ssh_url"`
	CloneURL        string      `json:"clone_url"`
	SvnURL          string      `json:"svn_url"`
	Homepage        string      `json:"homepage"`
	StargazersCount int         `json:"stargazers_count"`
	WatchersCount   int         `json:"watchers_count"`
	Language        string      `json:"language"`
	HasIssues       bool        `json:"has_issues"`
	HasProjects     bool        `json:"has_projects"`
	HasDownloads    bool        `json:"has_downloads"`
	HasWiki         bool        `json:"has_wiki"`
	HasPages        bool        `json:"has_pages"`
	HasDiscussions  bool        `json:"has_discussions"`
	ForksCount      int         `json:"forks_count"`
	MirrorURL       string      `json:"mirror_url"`
	Archived        bool        `json:"archived"`
	Disabled        bool        `json:"disabled"`
	OpenIssuesCount int         `json:"open_issues_count"`
	LicenseKey      string      `json:"license.key"`
	LicenseName     string      `json:"license.name"`
	AllowForking    bool        `json:"allow_forking"`
	IsTemplate      bool        `json:"is_template"`
	Topics          StringArray `json:"topics"`
	Visibility      string      `json:"visibility"`
	Forks           int         `json:"forks"`
	OpenIssues      int         `json:"open_issues"`
	Watchers        int         `json:"watchers"`
	DefaultBranch   string      `json:"default_branch"`
}
