// Package mockservice provides testify mock implementations of all service
// interfaces, used in handler-layer unit tests.
package mockservice

import (
	"context"
	"mime/multipart"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// AuthService is a mock implementation of service.AuthService.
type AuthService struct {
	mock.Mock
}

func (m *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

// JobService is a mock implementation of service.JobService.
type JobService struct {
	mock.Mock
}

func (m *JobService) CreateJob(ctx context.Context, recruiterID uuid.UUID, req dto.CreateJobRequest) (*dto.JobResponse, error) {
	args := m.Called(ctx, recruiterID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JobResponse), args.Error(1)
}

func (m *JobService) ListJobs(ctx context.Context, filter dto.JobFilter) ([]dto.JobResponse, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]dto.JobResponse), args.Get(1).(int64), args.Error(2)
}

func (m *JobService) GetMyJobs(ctx context.Context, recruiterID uuid.UUID, filter dto.RecruiterJobFilter) ([]dto.JobResponse, error) {
	args := m.Called(ctx, recruiterID, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.JobResponse), args.Error(1)
}

func (m *JobService) GetJobByID(ctx context.Context, id uuid.UUID) (*dto.JobResponse, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JobResponse), args.Error(1)
}

func (m *JobService) UpdateJob(ctx context.Context, recruiterID, jobID uuid.UUID, req dto.UpdateJobRequest) (*dto.JobResponse, error) {
	args := m.Called(ctx, recruiterID, jobID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.JobResponse), args.Error(1)
}

func (m *JobService) DeleteJob(ctx context.Context, recruiterID, jobID uuid.UUID) error {
	args := m.Called(ctx, recruiterID, jobID)
	return args.Error(0)
}

func (m *JobService) UpdateJobStages(ctx context.Context, recruiterID, jobID uuid.UUID, req dto.UpdateJobStagesRequest) ([]dto.StageResponse, error) {
	args := m.Called(ctx, recruiterID, jobID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.StageResponse), args.Error(1)
}

// ApplicationService is a mock implementation of service.ApplicationService.
type ApplicationService struct {
	mock.Mock
}

func (m *ApplicationService) Apply(ctx context.Context, candidateID, jobID uuid.UUID, coverLetter string, cvFile *multipart.FileHeader) (*dto.ApplicationResponse, error) {
	args := m.Called(ctx, candidateID, jobID, coverLetter, cvFile)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ApplicationResponse), args.Error(1)
}

func (m *ApplicationService) GetMyApplications(ctx context.Context, candidateID uuid.UUID) ([]dto.ApplicationResponse, error) {
	args := m.Called(ctx, candidateID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ApplicationResponse), args.Error(1)
}

func (m *ApplicationService) GetJobApplications(ctx context.Context, recruiterID, jobID uuid.UUID) ([]dto.ApplicationResponse, error) {
	args := m.Called(ctx, recruiterID, jobID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.ApplicationResponse), args.Error(1)
}

func (m *ApplicationService) Withdraw(ctx context.Context, candidateID, applicationID uuid.UUID) (*dto.ApplicationResponse, error) {
	args := m.Called(ctx, candidateID, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ApplicationResponse), args.Error(1)
}

func (m *ApplicationService) AdvanceStage(ctx context.Context, recruiterID, applicationID uuid.UUID) (*dto.ApplicationResponse, error) {
	args := m.Called(ctx, recruiterID, applicationID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ApplicationResponse), args.Error(1)
}

func (m *ApplicationService) UpdateStatus(ctx context.Context, recruiterID, applicationID uuid.UUID, status model.ApplicationStatus) (*dto.ApplicationResponse, error) {
	args := m.Called(ctx, recruiterID, applicationID, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ApplicationResponse), args.Error(1)
}
