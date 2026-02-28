package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"docker-workshop-assesment-grader/internal/models"
)

func Init(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := db.AutoMigrate(&models.Student{}); err != nil {
		return nil, fmt.Errorf("auto migrate database: %w", err)
	}

	return db, nil
}
