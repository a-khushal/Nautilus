package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/a-khushal/Nautilus/worker/models"
)

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&models.Job{}); err != nil {
		panic(err)
	}

	return db
}
