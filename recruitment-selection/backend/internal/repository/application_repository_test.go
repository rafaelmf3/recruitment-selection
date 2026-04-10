//go:build integration

package repository_test

import (
	"context"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/repository"
	"recruitment-selection/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedCandidateAndJob inserts a candidate user and an open job, returning both.
func seedCandidateAndJob(t *testing.T, db interface{ /* *gorm.DB */ }) (candidate *model.User, job *model.Job) {
	t.Helper()
	// This helper is intentionally empty - will be implemented in GREEN phase.
	// Tests below call seedRecruiter (defined in job_repository_test.go) directly.
	return nil, nil
}

func TestApplicationRepository_Create_And_FindByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	recruiter := seedRecruiter(t, userRepo)
	candidate := &model.User{
		ID: uuid.New(), Name: "Bob", Email: uuid.NewString() + "@gmail.com",
		PasswordHash: "h", Role: model.RoleCandidate,
	}
	require.NoError(t, userRepo.Create(context.Background(), candidate))

	job := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Dev", Description: "desc", Status: model.JobStatusOpen,
	}
	require.NoError(t, jobRepo.Create(context.Background(), job))

	app := &model.Application{
		ID:          uuid.New(),
		JobID:       job.ID,
		CandidateID: candidate.ID,
		Status:      model.ApplicationStatusPending,
		CoverLetter: "Motivated.",
	}
	require.NoError(t, appRepo.Create(context.Background(), app))

	found, err := appRepo.FindByID(context.Background(), app.ID)
	require.NoError(t, err)
	assert.Equal(t, app.ID, found.ID)
	assert.Equal(t, model.ApplicationStatusPending, found.Status)
}

func TestApplicationRepository_FindByID_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	appRepo := repository.NewApplicationRepository(db)

	_, err := appRepo.FindByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, apierror.ErrApplicationNotFound)
}

func TestApplicationRepository_FindByJobAndCandidate_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	appRepo := repository.NewApplicationRepository(db)

	_, err := appRepo.FindByJobAndCandidate(context.Background(), uuid.New(), uuid.New())
	assert.ErrorIs(t, err, apierror.ErrApplicationNotFound)
}

func TestApplicationRepository_UniqueConstraint(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	recruiter := seedRecruiter(t, userRepo)
	candidate := &model.User{
		ID: uuid.New(), Name: "Bob", Email: uuid.NewString() + "@gmail.com",
		PasswordHash: "h", Role: model.RoleCandidate,
	}
	require.NoError(t, userRepo.Create(context.Background(), candidate))

	job := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Dev", Description: "desc", Status: model.JobStatusOpen,
	}
	require.NoError(t, jobRepo.Create(context.Background(), job))

	app := &model.Application{
		ID: uuid.New(), JobID: job.ID, CandidateID: candidate.ID,
		Status: model.ApplicationStatusPending,
	}
	require.NoError(t, appRepo.Create(context.Background(), app))

	// Second application to same job by same candidate must fail.
	duplicate := &model.Application{
		ID: uuid.New(), JobID: job.ID, CandidateID: candidate.ID,
		Status: model.ApplicationStatusPending,
	}
	err := appRepo.Create(context.Background(), duplicate)
	assert.Error(t, err, "duplicate application must return an error")
}

func TestApplicationRepository_ListByCandidate(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	recruiter := seedRecruiter(t, userRepo)
	candidate := &model.User{
		ID: uuid.New(), Name: "Bob", Email: uuid.NewString() + "@gmail.com",
		PasswordHash: "h", Role: model.RoleCandidate,
	}
	require.NoError(t, userRepo.Create(context.Background(), candidate))

	// Create 2 jobs and apply to both.
	for i := 0; i < 2; i++ {
		job := &model.Job{
			ID: uuid.New(), RecruiterID: recruiter.ID,
			Title: "Job", Description: "desc", Status: model.JobStatusOpen,
		}
		require.NoError(t, jobRepo.Create(context.Background(), job))
		app := &model.Application{
			ID: uuid.New(), JobID: job.ID, CandidateID: candidate.ID,
			Status: model.ApplicationStatusPending,
		}
		require.NoError(t, appRepo.Create(context.Background(), app))
	}

	apps, err := appRepo.ListByCandidate(context.Background(), candidate.ID)
	require.NoError(t, err)
	assert.Len(t, apps, 2)
}

func TestApplicationRepository_Update_Status(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	recruiter := seedRecruiter(t, userRepo)
	candidate := &model.User{
		ID: uuid.New(), Name: "Bob", Email: uuid.NewString() + "@gmail.com",
		PasswordHash: "h", Role: model.RoleCandidate,
	}
	require.NoError(t, userRepo.Create(context.Background(), candidate))

	job := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Dev", Description: "desc", Status: model.JobStatusOpen,
	}
	require.NoError(t, jobRepo.Create(context.Background(), job))

	app := &model.Application{
		ID: uuid.New(), JobID: job.ID, CandidateID: candidate.ID,
		Status: model.ApplicationStatusPending,
	}
	require.NoError(t, appRepo.Create(context.Background(), app))

	app.Status = model.ApplicationStatusAccepted
	require.NoError(t, appRepo.Update(context.Background(), app))

	found, err := appRepo.FindByID(context.Background(), app.ID)
	require.NoError(t, err)
	assert.Equal(t, model.ApplicationStatusAccepted, found.Status)
}
