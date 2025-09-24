package processor

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/a-khushal/Nautilus/worker/models"
	"gorm.io/gorm"
)

func Process(jobID string, db *gorm.DB) {
	var job models.Job

	if err := db.First(&job, "id = ?", jobID).Error; err != nil {
		log.Println("Failed to fetch job from DB:", err)
		return
	}

	b, _ := json.Marshal(job)
	fmt.Println(string(b))
}
