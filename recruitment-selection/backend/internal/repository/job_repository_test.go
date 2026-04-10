//go:build integration

package repository_test

import (
	"context"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/repository"
	"recruitment-selection/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// seedRecruiter inserts a recruiter into the DB and returns it.
func seedRecruiter(t *testing.T, userRepo repository.UserRepository) *model.User {
	t.Helper()
	u := &model.User{
		ID:           uuid.New(),
		Name:         "Recruiter",
		Email:        uuid.NewString() + "@company.com",
		PasswordHash: "hashed",
		Role:         model.RoleRecruiter,
	}
	require.NoError(t, userRepo.Create(context.Background(), u))
	return u
}

func TestJobRepository_Create_And_FindByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	recruiter := seedRecruiter(t, userRepo)

	job := &model.Job{
		ID:          uuid.New(),
		RecruiterID: recruiter.ID,
		Title:       "Backend Engineer",
		Description: "Build APIs.",
		Status:      model.JobStatusOpen,
	}

	err := jobRepo.Create(context.Background(), job)
	require.NoError(t, err)

	found, err := jobRepo.FindByID(context.Background(), job.ID)
	require.NoError(t, err)
	assert.Equal(t, job.Title, found.Title)
	assert.Equal(t, model.JobStatusOpen, found.Status)
}

func TestJobRepository_FindByID_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	jobRepo := repository.NewJobRepository(db)

	_, err := jobRepo.FindByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, apierror.ErrJobNotFound)
}

func TestJobRepository_List_FilterByTitle(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	recruiter := seedRecruiter(t, userRepo)

	titles := []string{"Backend Engineer", "Frontend Developer", "DevOps Specialist"}
	for _, title := range titles {
		j := &model.Job{
			ID: uuid.New(), RecruiterID: recruiter.ID,
			Title: title, Description: "desc", Status: model.JobStatusOpen,
		}
		require.NoError(t, jobRepo.Create(context.Background(), j))
	}

	filter := dto.JobFilter{Title: "Backend", Page: 1, Limit: 20}
	jobs, total, err := jobRepo.List(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Backend Engineer", jobs[0].Title)
}

func TestJobRepository_List_FilterBySalary(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	recruiter := seedRecruiter(t, userRepo)

	lowMin, lowMax := 1000.0, 3000.0
	highMin, highMax := 8000.0, 12000.0

	lowJob := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Junior Dev", Description: "desc", Status: model.JobStatusOpen,
		SalaryMin: &lowMin, SalaryMax: &lowMax,
	}
	highJob := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Senior Dev", Description: "desc", Status: model.JobStatusOpen,
		SalaryMin: &highMin, SalaryMax: &highMax,
	}
	require.NoError(t, jobRepo.Create(context.Background(), lowJob))
	require.NoError(t, jobRepo.Create(context.Background(), highJob))

	// Filter: only jobs with salary_min >= 5000
	minFilter := 5000.0
	filter := dto.JobFilter{SalaryMin: &minFilter, Page: 1, Limit: 20}
	jobs, total, err := jobRepo.List(context.Background(), filter)

	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "Senior Dev", jobs[0].Title)
}

func TestJobRepository_Update(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	recruiter := seedRecruiter(t, userRepo)

	job := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Old Title", Description: "desc", Status: model.JobStatusOpen,
	}
	require.NoError(t, jobRepo.Create(context.Background(), job))

	job.Title = "New Title"
	job.Status = model.JobStatusClosed
	require.NoError(t, jobRepo.Update(context.Background(), job))

	found, err := jobRepo.FindByID(context.Background(), job.ID)
	require.NoError(t, err)
	assert.Equal(t, "New Title", found.Title)
	assert.Equal(t, model.JobStatusClosed, found.Status)
}

func TestJobRepository_Delete(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	recruiter := seedRecruiter(t, userRepo)

	job := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Temp Job", Description: "desc", Status: model.JobStatusOpen,
	}
	require.NoError(t, jobRepo.Create(context.Background(), job))
	require.NoError(t, jobRepo.Delete(context.Background(), job.ID))

	_, err := jobRepo.FindByID(context.Background(), job.ID)
	assert.ErrorIs(t, err, apierror.ErrJobNotFound)
}

func TestJobRepository_ReplaceStages(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	userRepo := repository.NewUserRepository(db)
	jobRepo := repository.NewJobRepository(db)
	recruiter := seedRecruiter(t, userRepo)

	job := &model.Job{
		ID: uuid.New(), RecruiterID: recruiter.ID,
		Title: "Backend Engineer", Description: "desc", Status: model.JobStatusOpen,
	}
	require.NoError(t, jobRepo.Create(context.Background(), job))

	newStages := []model.JobStage{
		{ID: uuid.New(), JobID: job.ID, Name: "Phone Screen", OrderIndex: 1},
		{ID: uuid.New(), JobID: job.ID, Name: "Take-home Test", OrderIndex: 2},
	}
	require.NoError(t, jobRepo.ReplaceStages(context.Background(), job.ID, newStages))

	stages, err := jobRepo.FindStagesByJobID(context.Background(), job.ID)
	require.NoError(t, err)
	assert.Len(t, stages, 2)
	assert.Equal(t, "Phone Screen", stages[0].Name)
}
