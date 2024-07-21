package models

import (
	"time"
)

type Config struct {
	ID        uint   `gorm:"primaryKey"`
	Key       string `gorm:"uniqueIndex"`
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
