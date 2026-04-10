package model

import (
	"time"

	"github.com/google/uuid"
)

// ApplicationStatus represents the overall decision on a candidate's application.
type ApplicationStatus string

const (
	ApplicationStatusPending    ApplicationStatus = "pending"
	ApplicationStatusInProgress ApplicationStatus = "in_progress"
	ApplicationStatusAccepted   ApplicationStatus = "accepted"
	ApplicationStatusRejected   ApplicationStatus = "rejected"
	ApplicationStatusWithdrawn  ApplicationStatus = "withdrawn"
)

// Application is the GORM model for the applications table.
// A candidate can apply to a job only once (unique constraint on job_id + candidate_id).
type Application struct {
	ID             uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	JobID          uuid.UUID         `gorm:"type:uuid;not null"`
	Job            Job               `gorm:"foreignKey:JobID"`
	CandidateID    uuid.UUID         `gorm:"type:uuid;not null"`
	Candidate      User              `gorm:"foreignKey:CandidateID"`
	CurrentStageID *uuid.UUID        `gorm:"type:uuid"`
	CurrentStage   *JobStage         `gorm:"foreignKey:CurrentStageID"`
	Status         ApplicationStatus `gorm:"type:application_status;not null;default:pending"`
	CoverLetter    string
	CVUrl          string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
