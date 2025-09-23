package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JobRequest struct {
	Type    string                 `json:"type" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"`
}

type Job struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type      string
	Payload   map[string]interface{} `gorm:"type:jsonb"`
	Status    string                 `gorm:"default:'pending'"`
	CreatedAt time.Time              `gorm:"autoCreateTime"`
}

var (
	ctx     = context.Background()
	JOBS    = "JOBS"
	DB_CONN = "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	db, err := gorm.Open(postgres.Open(DB_CONN), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&Job{}); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/jobs", func(c *gin.Context) {
		var req JobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		job := Job{
			Type:    req.Type,
			Payload: req.Payload,
		}

		if err := db.Create(&job).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to insert into DB"})
			return
		}

		jobJSON, err := json.Marshal(job)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to encode job"})
			return
		}

		if err := rdb.LPush(ctx, JOBS, jobJSON).Err(); err != nil {
			c.JSON(500, gin.H{"error": "Failed to push job"})
			return
		}

		c.JSON(200, gin.H{
			"message": "JOB_RECEIVED",
			"id":      job.ID,
		})
	})

	r.Run(":3000")
}
