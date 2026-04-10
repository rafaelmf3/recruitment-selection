// Package service defines the interfaces for all business logic operations.
// Handler tests use mock implementations of these interfaces.
package service

import (
	"context"
	"mime/multipart"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
)

// AuthService handles user registration and authentication.
type AuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
}

// JobService handles all job posting and pipeline stage operations.
type JobService interface {
	CreateJob(ctx context.Context, recruiterID uuid.UUID, req dto.CreateJobRequest) (*dto.JobResponse, error)
	ListJobs(ctx context.Context, filter dto.JobFilter) ([]dto.JobResponse, int64, error)
	GetMyJobs(ctx context.Context, recruiterID uuid.UUID, filter dto.RecruiterJobFilter) ([]dto.JobResponse, error)
	GetJobByID(ctx context.Context, id uuid.UUID) (*dto.JobResponse, error)
	UpdateJob(ctx context.Context, recruiterID, jobID uuid.UUID, req dto.UpdateJobRequest) (*dto.JobResponse, error)
	DeleteJob(ctx context.Context, recruiterID, jobID uuid.UUID) error
	UpdateJobStages(ctx context.Context, recruiterID, jobID uuid.UUID, req dto.UpdateJobStagesRequest) ([]dto.StageResponse, error)
}

// ApplicationService handles candidate applications and recruiter pipeline management.
type ApplicationService interface {
	Apply(ctx context.Context, candidateID, jobID uuid.UUID, coverLetter string, cvFile *multipart.FileHeader) (*dto.ApplicationResponse, error)
	GetMyApplications(ctx context.Context, candidateID uuid.UUID) ([]dto.ApplicationResponse, error)
	GetJobApplications(ctx context.Context, recruiterID, jobID uuid.UUID) ([]dto.ApplicationResponse, error)
	Withdraw(ctx context.Context, candidateID, applicationID uuid.UUID) (*dto.ApplicationResponse, error)
	AdvanceStage(ctx context.Context, recruiterID, applicationID uuid.UUID) (*dto.ApplicationResponse, error)
	UpdateStatus(ctx context.Context, recruiterID, applicationID uuid.UUID, status model.ApplicationStatus) (*dto.ApplicationResponse, error)
}
