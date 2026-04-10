package model

import (
	"time"

	"github.com/google/uuid"
)

// JobStatus represents whether a job posting is accepting applications.
type JobStatus string

const (
	JobStatusOpen      JobStatus = "open"
	JobStatusPaused    JobStatus = "paused"
	JobStatusClosed    JobStatus = "closed"
	JobStatusCancelled JobStatus = "cancelled"
)

// JobStatusTerminal returns true when no further transitions are allowed.
func JobStatusTerminal(s JobStatus) bool {
	return s == JobStatusClosed || s == JobStatusCancelled
}

// JobStatusAcceptingApplications returns true only when candidates can apply.
func JobStatusAcceptingApplications(s JobStatus) bool {
	return s == JobStatusOpen
}

// Job is the GORM model for the jobs table.
type Job struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	RecruiterID  uuid.UUID  `gorm:"type:uuid;not null"`
	Recruiter    User       `gorm:"foreignKey:RecruiterID"`
	Company      string
	Title        string     `gorm:"not null"`
	Description  string     `gorm:"not null"`
	Requirements string
	Location     string
	SalaryMin    *float64
	SalaryMax    *float64
	Status       JobStatus  `gorm:"type:job_status;not null;default:open"`
	Stages       []JobStage `gorm:"foreignKey:JobID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
