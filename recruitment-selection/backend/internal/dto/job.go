package dto

import (
	"recruitment-selection/internal/model"
	"time"

	"github.com/google/uuid"
)

// CreateJobRequest is the payload for POST /api/v1/recruiter/jobs.
type CreateJobRequest struct {
	Company      string       `json:"company"`
	Title        string       `json:"title"        binding:"required,min=3,max=255"`
	Description  string       `json:"description"  binding:"required"`
	Requirements string       `json:"requirements"`
	Location     string       `json:"location"`
	SalaryMin    *float64     `json:"salary_min"`
	SalaryMax    *float64     `json:"salary_max"`
	Stages       []StageInput `json:"stages"`
}

// UpdateJobRequest is the payload for PUT /api/v1/jobs/:id.
type UpdateJobRequest struct {
	Company      string            `json:"company"`
	Title        string            `json:"title"        binding:"omitempty,min=3,max=255"`
	Description  string            `json:"description"`
	Requirements string            `json:"requirements"`
	Location     string            `json:"location"`
	SalaryMin    *float64          `json:"salary_min"`
	SalaryMax    *float64          `json:"salary_max"`
	Status       model.JobStatus   `json:"status"       binding:"omitempty,oneof=open paused closed cancelled"`
}

// JobFilter contains query parameters for GET /api/v1/jobs (public listing).
type JobFilter struct {
	Q         string   `form:"q"`
	Location  string   `form:"location"`
	SalaryMin *float64 `form:"salary_min"`
	SalaryMax *float64 `form:"salary_max"`
	Page      int      `form:"page,default=1"`
	Limit     int      `form:"limit,default=20"`
}

// RecruiterJobFilter contains query parameters for GET /api/v1/recruiter/jobs.
type RecruiterJobFilter struct {
	Q       string `form:"q"`       // searches title and company (case-insensitive)
	Status  string `form:"status"`  // optional: open | paused | closed | cancelled
}

// StageInput is a single stage definition when updating a job's pipeline.
type StageInput struct {
	Name       string `json:"name"        binding:"required,min=1,max=100"`
	OrderIndex int    `json:"order_index" binding:"required,min=1"`
}

// UpdateJobStagesRequest is the payload for PUT /api/v1/jobs/:id/stages.
type UpdateJobStagesRequest struct {
	Stages []StageInput `json:"stages" binding:"required,min=1,dive"`
}

// StageResponse is the public representation of a job stage.
type StageResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	OrderIndex int       `json:"order_index"`
}

// JobResponse is the public representation of a job posting.
type JobResponse struct {
	ID           uuid.UUID       `json:"id"`
	Company      string          `json:"company"`
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	Requirements string          `json:"requirements"`
	Location     string          `json:"location"`
	SalaryMin    *float64        `json:"salary_min"`
	SalaryMax    *float64        `json:"salary_max"`
	Status       model.JobStatus `json:"status"`
	Recruiter    UserResponse    `json:"recruiter"`
	Stages       []StageResponse `json:"stages,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
}
