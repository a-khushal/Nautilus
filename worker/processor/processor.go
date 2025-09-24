package processor

import (
	"log"

	"github.com/a-khushal/Nautilus/worker/jobs"
	"github.com/a-khushal/Nautilus/worker/models"
	"gorm.io/gorm"
)

var (
	RUN_CODE = "RUN_CODE"
)

func Process(jobID string, db *gorm.DB) {
	var job models.Job

	if err := db.First(&job, "id = ?", jobID).Error; err != nil {
		log.Println("Failed to fetch job from DB:", err)
		return
	}

	switch job.Type {
	case RUN_CODE:
		jobs.RunCodeJob(job)
	default:
		log.Println("Unknown job type:", job.Type)
	}
}
