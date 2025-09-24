package routes

import (
	"context"
	"net/http"

	"github.com/a-khushal/Nautilus/api/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ctx  = context.Background()
	JOBS = "JOBS"
)

func RegisterJobRoutes(r *gin.Engine, db *gorm.DB, rdb *redis.Client) {
	r.POST("/jobs", func(c *gin.Context) {
		var req models.JobRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		job := models.Job{
			Type:    req.Type,
			Payload: req.Payload,
		}

		if err := db.Create(&job).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert into DB"})
			return
		}

		if err := rdb.LPush(ctx, JOBS, job.ID).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to push job"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "JOB_RECEIVED",
			"id":      job.ID,
		})
	})
}
