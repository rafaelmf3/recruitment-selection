// Package testutil provides shared helpers for all test layers.
package testutil

import (
	"recruitment-selection/internal/model"
	"time"

	"github.com/google/uuid"
)

// Fixed UUIDs used across tests so expectations can be set precisely.
var (
	RecruiterID  = uuid.MustParse("00000000-0000-0000-0000-000000000001")
	CandidateID  = uuid.MustParse("00000000-0000-0000-0000-000000000002")
	JobID        = uuid.MustParse("00000000-0000-0000-0000-000000000010")
	StageID1     = uuid.MustParse("00000000-0000-0000-0000-000000000020")
	StageID2     = uuid.MustParse("00000000-0000-0000-0000-000000000021")
	ApplicationID = uuid.MustParse("00000000-0000-0000-0000-000000000030")
)

// SalaryMin and SalaryMax are reusable salary pointers.
var (
	SalaryMin = ptr(3000.0)
	SalaryMax = ptr(6000.0)
)

func ptr[T any](v T) *T { return &v }

// NewRecruiter returns a recruiter User fixture.
func NewRecruiter() *model.User {
	return &model.User{
		ID:           RecruiterID,
		Name:         "Alice Recruiter",
		Email:        "alice@company.com",
		PasswordHash: "$2a$10$hashedpassword",
		Role:         model.RoleRecruiter,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// NewCandidate returns a candidate User fixture.
func NewCandidate() *model.User {
	return &model.User{
		ID:           CandidateID,
		Name:         "Bob Candidate",
		Email:        "bob@gmail.com",
		PasswordHash: "$2a$10$hashedpassword",
		Role:         model.RoleCandidate,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// NewJob returns an open Job fixture owned by the default recruiter.
func NewJob() *model.Job {
	return &model.Job{
		ID:          JobID,
		RecruiterID: RecruiterID,
		Title:       "Backend Engineer",
		Description: "Build cool things.",
		Location:    "Remote",
		SalaryMin:   SalaryMin,
		SalaryMax:   SalaryMax,
		Status:      model.JobStatusOpen,
		Stages:      DefaultStages(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// DefaultStages returns the 2-stage slice used in most tests.
func DefaultStages() []model.JobStage {
	return []model.JobStage{
		{ID: StageID1, JobID: JobID, Name: "Screening", OrderIndex: 1},
		{ID: StageID2, JobID: JobID, Name: "Technical Challenge", OrderIndex: 2},
	}
}

// NewApplication returns a pending Application fixture.
func NewApplication() *model.Application {
	return &model.Application{
		ID:          ApplicationID,
		JobID:       JobID,
		CandidateID: CandidateID,
		Status:      model.ApplicationStatusPending,
		CoverLetter: "I am very motivated.",
		CVUrl:       "uploads/bob_cv.pdf",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
