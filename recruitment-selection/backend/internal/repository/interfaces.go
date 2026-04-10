// Package repository defines the interfaces that all database access
// implementations must satisfy. This allows service layer tests to use
// mock implementations without touching the real database.
package repository

import (
	"context"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
)

// UserRepository handles persistence for users.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

// JobRepository handles persistence for jobs and their stages.
type JobRepository interface {
	Create(ctx context.Context, job *model.Job) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Job, error)
	FindByRecruiter(ctx context.Context, recruiterID uuid.UUID) ([]model.Job, error)
	List(ctx context.Context, filter dto.JobFilter) ([]model.Job, int64, error)
	Update(ctx context.Context, job *model.Job) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Stage management
	ReplaceStages(ctx context.Context, jobID uuid.UUID, stages []model.JobStage) error
	FindStageByID(ctx context.Context, id uuid.UUID) (*model.JobStage, error)
	FindStagesByJobID(ctx context.Context, jobID uuid.UUID) ([]model.JobStage, error)
}

// ApplicationRepository handles persistence for applications.
type ApplicationRepository interface {
	Create(ctx context.Context, app *model.Application) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Application, error)
	FindByJobAndCandidate(ctx context.Context, jobID, candidateID uuid.UUID) (*model.Application, error)
	ListByCandidate(ctx context.Context, candidateID uuid.UUID) ([]model.Application, error)
	ListByJob(ctx context.Context, jobID uuid.UUID) ([]model.Application, error)
	Update(ctx context.Context, app *model.Application) error
}
