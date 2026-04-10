package model

import (
	"time"

	"github.com/google/uuid"
)

// JobStage represents a single step in a job's selection pipeline.
// Stages are ordered by OrderIndex and belong to exactly one job.
type JobStage struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	JobID      uuid.UUID `gorm:"type:uuid;not null"`
	Name       string    `gorm:"not null"`
	OrderIndex int       `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// DefaultStageNames is the ordered list of stages created automatically
// when a recruiter creates a new job.
var DefaultStageNames = []string{
	"Screening",
	"Technical Challenge",
	"Team Interview",
	"Manager Interview",
	"Offer",
	"Hired",
}
