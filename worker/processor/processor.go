package processor

import (
	"context"
	"fmt"
	"log"

	"github.com/a-khushal/Nautilus/worker/jobs"
	"github.com/a-khushal/Nautilus/worker/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	RUN_CODE         = "RUN_CODE"
	JOB_COMPLETED    = "JOB_COMPLETED"
	UNKNOWN_JOB_TYPE = "UNKNOWN_JOB_TYPE"
	JOB_RESULTS      = "JOB_RESULTS"
)

func Process(ctx context.Context, jobID string, db *gorm.DB, rdb *redis.Client) {
	var job models.Job

	if err := db.First(&job, "id = ?", jobID).Error; err != nil {
		log.Println("Failed to fetch job from DB:", err)
		return
	}

	var result string

	switch job.Type {
	case RUN_CODE:
		res := jobs.RunCodeJob(job)
		result = string(res)
		fmt.Println("Job output:\n", result)
	default:
		log.Println("Unknown job type:", job.Type)
		result = "UNKNOWN_JOB_TYPE"
	}

	jobChannel := JOB_RESULTS + "_" + jobID
	if err := rdb.Publish(ctx, jobChannel, result).Err(); err != nil {
		log.Println("Failed to publish job result:", err)
		return
	}
	log.Println("Published result for job:", jobID)
}
