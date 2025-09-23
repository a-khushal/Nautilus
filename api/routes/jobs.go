package routes

import (
	"context"
	"encoding/json"

	"github.com/a-khushal/Nautilus/api/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var ctx = context.Background()

const JOBS = "JOBS"

func RegisterJobRoutes(r *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	r.POST("/jobs", func(c *gin.Context) {
		var req models.JobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		job := models.Job{
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
}
