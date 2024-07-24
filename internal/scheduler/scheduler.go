package scheduler

import (
	"encoding/json"
	"githubECS/rabbitmq"
	"log"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

func StartCommitMonitor(db *gorm.DB, rabbitCh *amqp.Channel) {
	s := gocron.NewScheduler(time.UTC)

	s.Every(1).Hours().Do(func() {
		var repos []string
		db.Raw("SELECT full_name FROM repositories").Scan(&repos)

		for _, repo := range repos {
			notifyCommitManager(repo, rabbitCh)
		}
	})

	s.StartBlocking()
}

func notifyCommitManager(repo string, rabbitCh *amqp.Channel) {
	queueName := "commit_manager_queue"
	body, err := json.Marshal(map[string]string{"repo": repo})
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	err = rabbitmq.PublishMessage(rabbitCh, queueName, body)
	if err != nil {
		log.Printf("Error publishing message to RabbitMQ: %v", err)
	}
}
