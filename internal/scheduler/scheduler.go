package scheduler

import (
	"github.com/go-co-op/gocron"
	"githubECS/internal/commit"
	"gorm.io/gorm"
	"time"
)

func StartScheduler(db *gorm.DB) {
	scheduler := gocron.NewScheduler(time.UTC)

	scheduler.Every(60).Minutes().Do(func() { commit.WatchCommits(db) })
	scheduler.StartAsync()
}
