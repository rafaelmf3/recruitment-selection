package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/repository"
)

// applicationService is the concrete implementation of ApplicationService.
type applicationService struct {
	appRepo   repository.ApplicationRepository
	jobRepo   repository.JobRepository
	uploadDir string
}

// NewApplicationService returns a new ApplicationService.
func NewApplicationService(
	appRepo repository.ApplicationRepository,
	jobRepo repository.JobRepository,
	uploadDir string,
) ApplicationService {
	return &applicationService{appRepo: appRepo, jobRepo: jobRepo, uploadDir: uploadDir}
}

func (s *applicationService) Apply(ctx context.Context, candidateID, jobID uuid.UUID, coverLetter string, cvFile *multipart.FileHeader) (*dto.ApplicationResponse, error) {
	// Verify the job exists and is open
	job, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	if !model.JobStatusAcceptingApplications(job.Status) {
		return nil, apierror.ErrJobNotAccepting
	}

	// Prevent duplicate applications
	_, err = s.appRepo.FindByJobAndCandidate(ctx, jobID, candidateID)
	if err == nil {
		return nil, apierror.ErrAlreadyApplied
	}

	// Save CV file if provided
	cvURL := ""
	if cvFile != nil {
		cvURL, err = s.saveCV(cvFile, candidateID)
		if err != nil {
			return nil, fmt.Errorf("apply: save cv: %w", err)
		}
	}

	app := &model.Application{
		ID:          uuid.New(),
		JobID:       jobID,
		CandidateID: candidateID,
		Status:      model.ApplicationStatusPending,
		CoverLetter: coverLetter,
		CVUrl:       cvURL,
	}

	if err := s.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("apply: create application: %w", err)
	}

	return toApplicationResponse(app), nil
}

func (s *applicationService) GetMyApplications(ctx context.Context, candidateID uuid.UUID) ([]dto.ApplicationResponse, error) {
	apps, err := s.appRepo.ListByCandidate(ctx, candidateID)
	if err != nil {
		return nil, fmt.Errorf("get my applications: %w", err)
	}

	responses := make([]dto.ApplicationResponse, len(apps))
	for i, a := range apps {
		responses[i] = *toApplicationResponse(&a)
	}
	return responses, nil
}

func (s *applicationService) GetJobApplications(ctx context.Context, recruiterID, jobID uuid.UUID) ([]dto.ApplicationResponse, error) {
	job, err := s.jobRepo.FindByID(ctx, jobID)
	if err != nil {
		return nil, err
	}

	if job.RecruiterID != recruiterID {
		return nil, apierror.ErrNotOwner
	}

	apps, err := s.appRepo.ListByJob(ctx, jobID)
	if err != nil {
		return nil, fmt.Errorf("get job applications: %w", err)
	}

	responses := make([]dto.ApplicationResponse, len(apps))
	for i, a := range apps {
		responses[i] = *toApplicationResponse(&a)
	}
	return responses, nil
}

func (s *applicationService) AdvanceStage(ctx context.Context, recruiterID, applicationID uuid.UUID) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app.Job.RecruiterID != recruiterID {
		return nil, apierror.ErrNotOwner
	}

	stages, err := s.jobRepo.FindStagesByJobID(ctx, app.JobID)
	if err != nil {
		return nil, fmt.Errorf("advance stage: fetch stages: %w", err)
	}

	nextStage, err := findNextStage(stages, app.CurrentStageID)
	if err != nil {
		return nil, err
	}

	app.CurrentStageID = &nextStage.ID
	app.CurrentStage = nextStage
	if app.Status == model.ApplicationStatusPending {
		app.Status = model.ApplicationStatusInProgress
	}

	if err := s.appRepo.Update(ctx, app); err != nil {
		return nil, fmt.Errorf("advance stage: update: %w", err)
	}

	return toApplicationResponse(app), nil
}

func (s *applicationService) Withdraw(ctx context.Context, candidateID, applicationID uuid.UUID) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app.CandidateID != candidateID {
		return nil, apierror.ErrNotOwner
	}

	if app.Status == model.ApplicationStatusWithdrawn {
		return toApplicationResponse(app), nil
	}

	app.Status = model.ApplicationStatusWithdrawn

	if err := s.appRepo.Update(ctx, app); err != nil {
		return nil, fmt.Errorf("withdraw: update: %w", err)
	}

	return toApplicationResponse(app), nil
}

func (s *applicationService) UpdateStatus(ctx context.Context, recruiterID, applicationID uuid.UUID, status model.ApplicationStatus) (*dto.ApplicationResponse, error) {
	app, err := s.appRepo.FindByID(ctx, applicationID)
	if err != nil {
		return nil, err
	}

	if app.Job.RecruiterID != recruiterID {
		return nil, apierror.ErrNotOwner
	}

	app.Status = status

	if err := s.appRepo.Update(ctx, app); err != nil {
		return nil, fmt.Errorf("update status: %w", err)
	}

	return toApplicationResponse(app), nil
}

// ---- Helpers ----------------------------------------------------------------

// findNextStage returns the stage that comes after currentStageID in the ordered list.
// If currentStageID is nil, it returns the first stage.
func findNextStage(stages []model.JobStage, currentStageID *uuid.UUID) (*model.JobStage, error) {
	if len(stages) == 0 {
		return nil, apierror.ErrStageNotFound
	}

	// No current stage - advance to first
	if currentStageID == nil {
		return &stages[0], nil
	}

	for i, s := range stages {
		if s.ID == *currentStageID {
			if i+1 >= len(stages) {
				return nil, apierror.ErrNoNextStage
			}
			return &stages[i+1], nil
		}
	}

	return nil, apierror.ErrStageNotFound
}

// saveCV persists the uploaded file to the upload directory and returns its path.
func (s *applicationService) saveCV(file *multipart.FileHeader, candidateID uuid.UUID) (string, error) {
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%s%s", candidateID.String(), uuid.New().String(), ext)
	dst := filepath.Join(s.uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := out.ReadFrom(src); err != nil {
		return "", err
	}

	return dst, nil
}

func toApplicationResponse(a *model.Application) *dto.ApplicationResponse {
	resp := &dto.ApplicationResponse{
		ID:          a.ID,
		JobID:       a.JobID,
		CandidateID: a.CandidateID,
		Status:      a.Status,
		CoverLetter: a.CoverLetter,
		CVUrl:       a.CVUrl,
		CreatedAt:   a.CreatedAt,
	}

	// Populate job when preloaded
	if a.Job.ID != (uuid.UUID{}) {
		job := toJobResponse(&a.Job)
		resp.Job = job
	}

	// Populate candidate when preloaded
	if a.Candidate.ID != (uuid.UUID{}) {
		resp.Candidate = toUserResponse(&a.Candidate)
	}

	if a.CurrentStageID != nil {
		resp.CurrentStageID = a.CurrentStageID
	}

	if a.CurrentStage != nil {
		stage := toStageResponse(a.CurrentStage)
		resp.CurrentStage = &stage
	}

	return resp
}
