package models

import (
	"time"

	"gorm.io/datatypes"
)

type JobRequest struct {
	Type    string                 `json:"type" binding:"required"`
	Payload map[string]interface{} `json:"payload" binding:"required"`
}

type Job struct {
	ID        string `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Type      string
	Payload   datatypes.JSONMap `gorm:"type:jsonb"`
	Status    string            `gorm:"default:'pending'"`
	CreatedAt time.Time         `gorm:"autoCreateTime"`
}
