package dto

import (
	"recruitment-selection/internal/model"
	"time"

	"github.com/google/uuid"
)

// UpdateApplicationStatusRequest is the payload for PATCH /api/v1/applications/:id/status.
// Only recruiter-controlled statuses are allowed here; "withdrawn" is handled
// by the candidate's own PATCH /applications/:id/withdraw endpoint.
type UpdateApplicationStatusRequest struct {
	Status model.ApplicationStatus `json:"status" binding:"required,oneof=accepted rejected"`
}

// ApplicationResponse is the public representation of an application.
type ApplicationResponse struct {
	ID             uuid.UUID                `json:"id"`
	JobID          uuid.UUID                `json:"job_id"`
	Job            *JobResponse             `json:"job,omitempty"`
	CandidateID    uuid.UUID                `json:"candidate_id"`
	Candidate      *UserResponse            `json:"candidate,omitempty"`
	CurrentStageID *uuid.UUID               `json:"current_stage_id,omitempty"`
	CurrentStage   *StageResponse           `json:"current_stage,omitempty"`
	Status         model.ApplicationStatus  `json:"status"`
	CoverLetter    string                   `json:"cover_letter"`
	CVUrl          string                   `json:"cv_url"`
	CreatedAt      time.Time                `json:"created_at"`
}
