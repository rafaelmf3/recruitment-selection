// Package mockrepo provides testify mock implementations of all repository
// interfaces, used in service-layer unit tests.
package mockrepo

import (
	"context"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// UserRepository is a mock implementation of repository.UserRepository.
type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

// JobRepository is a mock implementation of repository.JobRepository.
type JobRepository struct {
	mock.Mock
}

func (m *JobRepository) Create(ctx context.Context, job *model.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *JobRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Job, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Job), args.Error(1)
}

func (m *JobRepository) FindByRecruiter(ctx context.Context, recruiterID uuid.UUID) ([]model.Job, error) {
	args := m.Called(ctx, recruiterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Job), args.Error(1)
}

func (m *JobRepository) List(ctx context.Context, filter dto.JobFilter) ([]model.Job, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]model.Job), args.Get(1).(int64), args.Error(2)
}

func (m *JobRepository) Update(ctx context.Context, job *model.Job) error {
	args := m.Called(ctx, job)
	return args.Error(0)
}

func (m *JobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *JobRepository) ReplaceStages(ctx context.Context, jobID uuid.UUID, stages []model.JobStage) error {
	args := m.Called(ctx, jobID, stages)
	return args.Error(0)
}

func (m *JobRepository) FindStageByID(ctx context.Context, id uuid.UUID) (*model.JobStage, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.JobStage), args.Error(1)
}

func (m *JobRepository) FindStagesByJobID(ctx context.Context, jobID uuid.UUID) ([]model.JobStage, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.JobStage), args.Error(1)
}

// ApplicationRepository is a mock implementation of repository.ApplicationRepository.
type ApplicationRepository struct {
	mock.Mock
}

func (m *ApplicationRepository) Create(ctx context.Context, app *model.Application) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}

func (m *ApplicationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Application), args.Error(1)
}

func (m *ApplicationRepository) FindByJobAndCandidate(ctx context.Context, jobID, candidateID uuid.UUID) (*model.Application, error) {
	args := m.Called(ctx, jobID, candidateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Application), args.Error(1)
}

func (m *ApplicationRepository) ListByCandidate(ctx context.Context, candidateID uuid.UUID) ([]model.Application, error) {
	args := m.Called(ctx, candidateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Application), args.Error(1)
}

func (m *ApplicationRepository) ListByJob(ctx context.Context, jobID uuid.UUID) ([]model.Application, error) {
	args := m.Called(ctx, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Application), args.Error(1)
}

func (m *ApplicationRepository) Update(ctx context.Context, app *model.Application) error {
	args := m.Called(ctx, app)
	return args.Error(0)
}
