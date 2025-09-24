package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/a-khushal/Nautilus/api/models"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ctx           = context.Background()
	JOBS          = "JOBS"
	JOB_RESULTS   = "JOB_RESULTS"
	JOB_SUBMITTED = "JOB_SUBMITTED"
	JOB_COMPLETED = "JOB_COMPLETED"
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

		jobChannel := JOB_RESULTS + "_" + job.ID
		pubsub := rdb.Subscribe(ctx, jobChannel)
		defer pubsub.Close()
		ch := pubsub.Channel()

		timeout := time.After(30 * time.Second)

		for {
			select {
			case msg := <-ch:
				c.JSON(http.StatusOK, gin.H{
					"message": JOB_COMPLETED,
					"id":      job.ID,
					"result":  msg.Payload,
				})
				return
			case <-timeout:
				c.JSON(http.StatusOK, gin.H{
					"message": JOB_SUBMITTED,
					"id":      job.ID,
				})
				return
			}
		}
	})
}
