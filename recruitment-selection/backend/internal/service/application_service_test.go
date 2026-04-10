package service_test

import (
	"context"
	"errors"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/model"
	mockrepo "recruitment-selection/internal/mock/repository"
	"recruitment-selection/internal/service"
	"recruitment-selection/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newApplicationService is intentionally not used as a shared helper because
// each test needs t.TempDir() scoped to its own *testing.T.
// Call service.NewApplicationService(appRepo, jobRepo, t.TempDir()) directly.

// ---- Apply ------------------------------------------------------------------

func TestApplicationService_Apply_CandidateSuccess(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	job := testutil.NewJob()
	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)
	appRepo.On("FindByJobAndCandidate", context.Background(), testutil.JobID, testutil.CandidateID).
		Return(nil, apierror.ErrApplicationNotFound)
	appRepo.On("Create", context.Background(), testutil.AnyApplication()).
		Return(nil)

	resp, err := svc.Apply(context.Background(), testutil.CandidateID, testutil.JobID, "Cover letter text", nil)

	require.NoError(t, err)
	assert.Equal(t, testutil.JobID, resp.JobID)
	assert.Equal(t, testutil.CandidateID, resp.CandidateID)
	assert.Equal(t, model.ApplicationStatusPending, resp.Status)
	appRepo.AssertExpectations(t)
}

func TestApplicationService_Apply_JobClosed(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	closedJob := testutil.NewJob()
	closedJob.Status = model.JobStatusClosed

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(closedJob, nil)

	resp, err := svc.Apply(context.Background(), testutil.CandidateID, testutil.JobID, "Cover letter", nil)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrJobNotAccepting))
	appRepo.AssertNotCalled(t, "Create")
}

func TestApplicationService_Apply_PausedJob(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	pausedJob := testutil.NewJob()
	pausedJob.Status = model.JobStatusPaused

	jobRepo.On("FindByID", context.Background(), testutil.JobID).Return(pausedJob, nil)

	resp, err := svc.Apply(context.Background(), testutil.CandidateID, testutil.JobID, "Cover letter", nil)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrJobNotAccepting))
	appRepo.AssertNotCalled(t, "Create")
}

func TestApplicationService_Apply_AlreadyApplied(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	job := testutil.NewJob()
	existing := testutil.NewApplication()

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)
	appRepo.On("FindByJobAndCandidate", context.Background(), testutil.JobID, testutil.CandidateID).
		Return(existing, nil)

	resp, err := svc.Apply(context.Background(), testutil.CandidateID, testutil.JobID, "Cover letter", nil)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrAlreadyApplied))
	appRepo.AssertNotCalled(t, "Create")
}

// ---- GetMyApplications ------------------------------------------------------

func TestApplicationService_GetMyApplications(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	apps := []model.Application{*testutil.NewApplication()}
	appRepo.On("ListByCandidate", context.Background(), testutil.CandidateID).
		Return(apps, nil)

	result, err := svc.GetMyApplications(context.Background(), testutil.CandidateID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	appRepo.AssertExpectations(t)
}

// ---- GetJobApplications -----------------------------------------------------

func TestApplicationService_GetJobApplications_OwnerSuccess(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	job := testutil.NewJob()
	apps := []model.Application{*testutil.NewApplication()}

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)
	appRepo.On("ListByJob", context.Background(), testutil.JobID).
		Return(apps, nil)

	result, err := svc.GetJobApplications(context.Background(), testutil.RecruiterID, testutil.JobID)

	require.NoError(t, err)
	assert.Len(t, result, 1)
	jobRepo.AssertExpectations(t)
}

func TestApplicationService_GetJobApplications_NotOwner(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	job := testutil.NewJob() // owned by RecruiterID
	anotherRecruiter := testutil.CandidateID

	jobRepo.On("FindByID", context.Background(), testutil.JobID).
		Return(job, nil)

	result, err := svc.GetJobApplications(context.Background(), anotherRecruiter, testutil.JobID)

	assert.Nil(t, result)
	assert.True(t, errors.Is(err, apierror.ErrNotOwner))
	appRepo.AssertNotCalled(t, "ListByJob")
}

// ---- AdvanceStage -----------------------------------------------------------

func TestApplicationService_AdvanceStage_Success(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	app := testutil.NewApplication()
	app.Job = *testutil.NewJob()
	// Application currently has no stage - will advance to first stage.

	stages := testutil.DefaultStages()

	appRepo.On("FindByID", context.Background(), testutil.ApplicationID).
		Return(app, nil)
	jobRepo.On("FindStagesByJobID", context.Background(), testutil.JobID).
		Return(stages, nil)
	appRepo.On("Update", context.Background(), testutil.AnyApplication()).
		Return(nil)

	resp, err := svc.AdvanceStage(context.Background(), testutil.RecruiterID, testutil.ApplicationID)

	require.NoError(t, err)
	require.NotNil(t, resp.CurrentStage)
	assert.Equal(t, "Screening", resp.CurrentStage.Name)
	appRepo.AssertExpectations(t)
}

func TestApplicationService_AdvanceStage_AlreadyAtLastStage(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	stages := testutil.DefaultStages()
	lastStageID := stages[len(stages)-1].ID

	app := testutil.NewApplication()
	app.Job = *testutil.NewJob()
	app.CurrentStageID = &lastStageID // already at the last stage

	appRepo.On("FindByID", context.Background(), testutil.ApplicationID).
		Return(app, nil)
	jobRepo.On("FindStagesByJobID", context.Background(), testutil.JobID).
		Return(stages, nil)

	resp, err := svc.AdvanceStage(context.Background(), testutil.RecruiterID, testutil.ApplicationID)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrNoNextStage))
	appRepo.AssertNotCalled(t, "Update")
}

// ---- UpdateStatus -----------------------------------------------------------

func TestApplicationService_UpdateStatus_Accept(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	app := testutil.NewApplication()
	app.Job = *testutil.NewJob()

	appRepo.On("FindByID", context.Background(), testutil.ApplicationID).
		Return(app, nil)
	appRepo.On("Update", context.Background(), testutil.AnyApplication()).
		Return(nil)

	resp, err := svc.UpdateStatus(context.Background(), testutil.RecruiterID, testutil.ApplicationID, model.ApplicationStatusAccepted)

	require.NoError(t, err)
	assert.Equal(t, model.ApplicationStatusAccepted, resp.Status)
	appRepo.AssertExpectations(t)
}

func TestApplicationService_UpdateStatus_NotOwner(t *testing.T) {
	appRepo := new(mockrepo.ApplicationRepository)
	jobRepo := new(mockrepo.JobRepository)
	svc := service.NewApplicationService(appRepo, jobRepo, t.TempDir())

	app := testutil.NewApplication()
	app.Job = *testutil.NewJob() // RecruiterID owns the job
	anotherRecruiter := testutil.CandidateID

	appRepo.On("FindByID", context.Background(), testutil.ApplicationID).
		Return(app, nil)

	resp, err := svc.UpdateStatus(context.Background(), anotherRecruiter, testutil.ApplicationID, model.ApplicationStatusRejected)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrNotOwner))
	appRepo.AssertNotCalled(t, "Update")
}
