package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/repository"
)

// jobService is the concrete implementation of JobService.
type jobService struct {
	jobRepo repository.JobRepository
}

// NewJobService returns a new JobService.
func NewJobService(jobRepo repository.JobRepository) JobService {
	return &jobService{jobRepo: jobRepo}
}

func (s *jobService) CreateJob(ctx context.Context, recruiterID uuid.UUID, req dto.CreateJobRequest) (*dto.JobResponse, error) {
	if req.SalaryMin != nil && req.SalaryMax != nil && *req.SalaryMin > *req.SalaryMax {
		return nil, apierror.ErrInvalidSalary
	}

	job := &model.Job{
		ID:           uuid.New(),
		RecruiterID:  recruiterID,
		Company:      req.Company,
		Title:        req.Title,
		Description:  req.Description,
		Requirements: req.Requirements,
		Location:     req.Location,
		SalaryMin:    req.SalaryMin,
		SalaryMax:    req.SalaryMax,
		Status:       model.JobStatusOpen,
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("create job: %w", err)
	}

	// Create initial stages if provided
	if len(req.Stages) > 0 {
		stages := make([]model.JobStage, len(req.Stages))
		for i, s := range req.Stages {
			stages[i] = model.JobStage{
				ID:         uuid.New(),
				JobID:      job.ID,
				Name:       s.Name,
				OrderIndex: s.OrderIndex,
			}
		}
		if err := s.jobRepo.ReplaceStages(ctx, job.ID, stages); err != nil {
			return nil, fmt.Errorf("create job stages: %w", err)
		}
		job.Stages = stages
	}

	return toJobResponse(job), nil
}

func (s *jobService) ListJobs(ctx context.Context, filter dto.JobFilter) ([]dto.JobResponse, int64, error) {
	jobs, total, err := s.jobRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("list jobs: %w", err)
	}

	responses := make([]dto.JobResponse, len(jobs))
	for i, j := range jobs {
		responses[i] = *toJobResponse(&j)
	}
	return responses, total, nil
}

func (s *jobService) GetMyJobs(ctx context.Context, recruiterID uuid.UUID) ([]dto.JobResponse, error) {
	jobs, err := s.jobRepo.FindByRecruiter(ctx, recruiterID)
	if err != nil {
		return nil, fmt.Errorf("get my jobs: %w", err)
	}

	responses := make([]dto.JobResponse, len(jobs))
	for i, j := range jobs {
		responses[i] = *toJobResponse(&j)
	}
	return responses, nil
}

func (s *jobService) GetJobByID(ctx context.Context, id uuid.UUID) (*dto.JobResponse, error) {
	job, err := s.jobRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toJobResponse(job), nil
}

func (s *jobService) UpdateJob(ctx context.Context, recruiterID, jobID uuid.UUID, req dto.UpdateJobRequest) (*dto.JobResponse, error) {
	job, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if job.RecruiterID != recruiterID {
		return nil, apierror.ErrNotOwner
	}

	if req.SalaryMin != nil && req.SalaryMax != nil && *req.SalaryMin > *req.SalaryMax {
		return nil, apierror.ErrInvalidSalary
	}

	// Apply partial updates - only overwrite non-zero values
	if req.Company != "" {
		job.Company = req.Company
	}
	if req.Title != "" {
		job.Title = req.Title
	}
	if req.Description != "" {
		job.Description = req.Description
	}
	if req.Requirements != "" {
		job.Requirements = req.Requirements
	}
	if req.Location != "" {
		job.Location = req.Location
	}
	if req.SalaryMin != nil {
		job.SalaryMin = req.SalaryMin
	}
	if req.SalaryMax != nil {
		job.SalaryMax = req.SalaryMax
	}
	if req.Status != "" {
		if err := validateStatusTransition(job.Status, req.Status); err != nil {
			return nil, err
		}
		job.Status = req.Status
	}

	if err := s.jobRepo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("update job: %w", err)
	}

	return toJobResponse(job), nil
}

func (s *jobService) DeleteJob(ctx context.Context, recruiterID, jobID uuid.UUID) error {
	job, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return err
	}

	if job.RecruiterID != recruiterID {
		return apierror.ErrNotOwner
	}

	return s.jobRepo.Delete(ctx, jobID)
}

func (s *jobService) UpdateJobStages(ctx context.Context, recruiterID, jobID uuid.UUID, req dto.UpdateJobStagesRequest) ([]dto.StageResponse, error) {
	job, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if job.RecruiterID != recruiterID {
		return nil, apierror.ErrNotOwner
	}

	stages := make([]model.JobStage, len(req.Stages))
	for i, s := range req.Stages {
		stages[i] = model.JobStage{
			ID:         uuid.New(),
			JobID:      jobID,
			Name:       s.Name,
			OrderIndex: s.OrderIndex,
		}
	}

	if err := s.jobRepo.ReplaceStages(ctx, jobID, stages); err != nil {
		return nil, fmt.Errorf("update stages: %w", err)
	}

	responses := make([]dto.StageResponse, len(stages))
	for i, st := range stages {
		responses[i] = toStageResponse(&st)
	}
	return responses, nil
}

// ---- Status transition guard ------------------------------------------------

// validTransitions defines all legal from->to status moves.
var validTransitions = map[model.JobStatus][]model.JobStatus{
	model.JobStatusOpen:   {model.JobStatusPaused, model.JobStatusClosed, model.JobStatusCancelled},
	model.JobStatusPaused: {model.JobStatusOpen, model.JobStatusClosed, model.JobStatusCancelled},
	// closed and cancelled are terminal - no outgoing transitions
}

func validateStatusTransition(from, to model.JobStatus) error {
	if from == to {
		return nil
	}
	allowed, ok := validTransitions[from]
	if !ok {
		return apierror.ErrInvalidTransition
	}
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return apierror.ErrInvalidTransition
}

// ---- Mapping helpers --------------------------------------------------------

func toJobResponse(j *model.Job) *dto.JobResponse {
	resp := &dto.JobResponse{
		ID:           j.ID,
		Company:      j.Company,
		Title:        j.Title,
		Description:  j.Description,
		Requirements: j.Requirements,
		Location:     j.Location,
		SalaryMin:    j.SalaryMin,
		SalaryMax:    j.SalaryMax,
		Status:       j.Status,
		CreatedAt:    j.CreatedAt,
		Recruiter: dto.UserResponse{
			ID:    j.Recruiter.ID.String(),
			Name:  j.Recruiter.Name,
			Email: j.Recruiter.Email,
			Role:  j.Recruiter.Role,
		},
	}

	for _, st := range j.Stages {
		resp.Stages = append(resp.Stages, toStageResponse(&st))
	}
	return resp
}

func toStageResponse(s *model.JobStage) dto.StageResponse {
	return dto.StageResponse{
		ID:         s.ID,
		Name:       s.Name,
		OrderIndex: s.OrderIndex,
	}
}
