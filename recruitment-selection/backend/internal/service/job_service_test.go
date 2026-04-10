package service_test

import (
	"context"
	"errors"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	mockrepo "recruitment-selection/internal/mock/repository"
	"recruitment-selection/internal/service"
	"recruitment-selection/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// newJobService creates a jobService with an ephemeral appRepo mock.
// Use newJobServiceWithAppRepo when you need to assert on appRepo behaviour.
func newJobService(jobRepo *mockrepo.JobRepository) service.JobService {
	return service.NewJobService(jobRepo, new(mockrepo.ApplicationRepository))
}

// newJobServiceWithAppRepo creates a jobService returning an explicit appRepo
// mock so tests can configure expectations on it.
func newJobServiceWithAppRepo(
	jobRepo *mockrepo.JobRepository,
	appRepo *mockrepo.ApplicationRepository,
) service.JobService {
	return service.NewJobService(jobRepo, appRepo)
}

// ---- CreateJob --------------------------------------------------------------

func TestJobService_CreateJob_RecruiterSuccess(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	req := dto.CreateJobRequest{
		Title:       "Backend Engineer",
		Description: "Build APIs.",
		SalaryMin:   testutil.SalaryMin,
		SalaryMax:   testutil.SalaryMax,
	}

	jobRepo.On("Create", context.Background(), testutil.AnyJob()).
		Return(nil)

	resp, err := svc.CreateJob(context.Background(), testutil.RecruiterID, req)

	require.NoError(t, err)
	assert.Equal(t, req.Title, resp.Title)
	assert.Equal(t, model.JobStatusOpen, resp.Status)
	jobRepo.AssertExpectations(t)
}

func TestJobService_CreateJob_InvalidSalaryRange(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	min := 9000.0
	max := 3000.0
	req := dto.CreateJobRequest{
		Title:       "Backend Engineer",
		Description: "Build APIs.",
		SalaryMin:   &min,
		SalaryMax:   &max,
	}

	resp, err := svc.CreateJob(context.Background(), testutil.RecruiterID, req)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrInvalidSalary))
	// Repository must not be called when validation fails.
	jobRepo.AssertNotCalled(t, "Create")
}

// ---- ListJobs ---------------------------------------------------------------

func TestJobService_ListJobs_ReturnsOpenJobs(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	jobs := []model.Job{*testutil.NewJob()}
	filter := dto.JobFilter{Page: 1, Limit: 20}

	jobRepo.On("List", context.Background(), filter).
		Return(jobs, int64(1), nil)

	result, total, err := svc.ListJobs(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, jobs[0].Title, result[0].Title)
	jobRepo.AssertExpectations(t)
}

func TestJobService_ListJobs_FilterBySalary(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	min := 5000.0
	filter := dto.JobFilter{SalaryMin: &min, Page: 1, Limit: 20}

	jobRepo.On("List", context.Background(), filter).
		Return([]model.Job{}, int64(0), nil)

	result, total, err := svc.ListJobs(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, int64(0), total)
	assert.Empty(t, result)
	jobRepo.AssertExpectations(t)
}

// ---- GetJobByID -------------------------------------------------------------

func TestJobService_GetJobByID_Found(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob()
	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)

	resp, err := svc.GetJobByID(context.Background(), testutil.JobID)

	require.NoError(t, err)
	assert.Equal(t, job.Title, resp.Title)
	assert.Len(t, resp.Stages, len(job.Stages))
	jobRepo.AssertExpectations(t)
}

func TestJobService_GetJobByID_NotFound(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(nil, apierror.ErrJobNotFound)

	resp, err := svc.GetJobByID(context.Background(), testutil.JobID)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrJobNotFound))
	jobRepo.AssertExpectations(t)
}

// ---- UpdateJob --------------------------------------------------------------

func TestJobService_UpdateJob_OwnerSuccess(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob()
	req := dto.UpdateJobRequest{Title: "Senior Backend Engineer"}

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)
	jobRepo.On("Update", context.Background(), testutil.AnyJob()).
		Return(nil)

	resp, err := svc.UpdateJob(context.Background(), testutil.RecruiterID, testutil.JobID, req)

	require.NoError(t, err)
	assert.Equal(t, req.Title, resp.Title)
	jobRepo.AssertExpectations(t)
}

func TestJobService_UpdateJob_NotOwner(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob() // owned by RecruiterID
	anotherRecruiterID := testutil.CandidateID

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)

	req := dto.UpdateJobRequest{Title: "New Title"}
	resp, err := svc.UpdateJob(context.Background(), anotherRecruiterID, testutil.JobID, req)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrNotOwner))
	jobRepo.AssertNotCalled(t, "Update")
}

// ---- DeleteJob --------------------------------------------------------------

func TestJobService_DeleteJob_OwnerSuccess(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob()
	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)
	jobRepo.On("Delete", context.Background(), testutil.JobID).
		Return(nil)

	err := svc.DeleteJob(context.Background(), testutil.RecruiterID, testutil.JobID)

	require.NoError(t, err)
	jobRepo.AssertExpectations(t)
}

func TestJobService_DeleteJob_NotOwner(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob()
	anotherRecruiterID := testutil.CandidateID

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)

	err := svc.DeleteJob(context.Background(), anotherRecruiterID, testutil.JobID)

	assert.True(t, errors.Is(err, apierror.ErrNotOwner))
	jobRepo.AssertNotCalled(t, "Delete")
}

// ---- Status transitions -----------------------------------------------------

func TestJobService_UpdateJob_PauseFromOpen(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob() // status: open
	req := dto.UpdateJobRequest{Status: model.JobStatusPaused}

	jobRepo.On("FindByID", context.Background(), testutil.JobID).Return(job, nil)
	jobRepo.On("Update", context.Background(), testutil.AnyJob()).Return(nil)

	resp, err := svc.UpdateJob(context.Background(), testutil.RecruiterID, testutil.JobID, req)

	require.NoError(t, err)
	assert.Equal(t, model.JobStatusPaused, resp.Status)
}

func TestJobService_UpdateJob_TerminalJobIsImmutable(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob()
	job.Status = model.JobStatusClosed // terminal state — no edits allowed
	req := dto.UpdateJobRequest{Title: "Any change should be rejected"}

	jobRepo.On("FindByID", context.Background(), testutil.JobID).Return(job, nil)

	resp, err := svc.UpdateJob(context.Background(), testutil.RecruiterID, testutil.JobID, req)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrJobTerminal))
	jobRepo.AssertNotCalled(t, "Update")
}

func TestJobService_UpdateJob_CloseJobRejectsActiveApplications(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	appRepo := new(mockrepo.ApplicationRepository)
	svc := newJobServiceWithAppRepo(jobRepo, appRepo)

	job := testutil.NewJob() // status: open
	req := dto.UpdateJobRequest{Status: model.JobStatusClosed}

	jobRepo.On("FindByID", context.Background(), testutil.JobID).Return(job, nil)
	appRepo.On("RejectActiveApplications", context.Background(), testutil.JobID).Return(nil)
	jobRepo.On("Update", context.Background(), testutil.AnyJob()).Return(nil)

	resp, err := svc.UpdateJob(context.Background(), testutil.RecruiterID, testutil.JobID, req)

	require.NoError(t, err)
	assert.Equal(t, model.JobStatusClosed, resp.Status)
	appRepo.AssertExpectations(t)
	jobRepo.AssertExpectations(t)
}

// ---- UpdateJobStages --------------------------------------------------------

func TestJobService_UpdateJobStages_ReplacesAll(t *testing.T) {
	jobRepo := new(mockrepo.JobRepository)
	svc := newJobService(jobRepo)

	job := testutil.NewJob()
	req := dto.UpdateJobStagesRequest{
		Stages: []dto.StageInput{
			{Name: "Phone Screen", OrderIndex: 1},
			{Name: "Coding Test", OrderIndex: 2},
			{Name: "Offer", OrderIndex: 3},
		},
	}

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)
	jobRepo.On("ReplaceStages", context.Background(), testutil.JobID, mock.AnythingOfType("[]model.JobStage")).
		Return(nil)

	stages, err := svc.UpdateJobStages(context.Background(), testutil.RecruiterID, testutil.JobID, req)

	require.NoError(t, err)
	assert.Len(t, stages, 3)
	assert.Equal(t, "Phone Screen", stages[0].Name)
	jobRepo.AssertExpectations(t)
}
