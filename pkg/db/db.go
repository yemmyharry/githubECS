package db

import (
	"githubECS/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&models.Repository{}, &models.Commit{}, &models.Config{})
	return db, nil
}
