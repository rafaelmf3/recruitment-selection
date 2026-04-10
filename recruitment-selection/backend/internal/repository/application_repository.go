package repository

import (
	"context"
	"errors"
	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// applicationRepository is the PostgreSQL implementation of ApplicationRepository.
type applicationRepository struct {
	db *gorm.DB
}

// NewApplicationRepository returns a new ApplicationRepository backed by the given DB connection.
func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(ctx context.Context, app *model.Application) error {
	return r.db.WithContext(ctx).Create(app).Error
}

func (r *applicationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	var app model.Application
	err := r.db.WithContext(ctx).
		Preload("Job.Recruiter").
		Preload("Candidate").
		Preload("CurrentStage").
		First(&app, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.ErrApplicationNotFound
	}
	return &app, err
}

func (r *applicationRepository) FindByJobAndCandidate(ctx context.Context, jobID, candidateID uuid.UUID) (*model.Application, error) {
	var app model.Application
	err := r.db.WithContext(ctx).
		Where("job_id = ? AND candidate_id = ?", jobID, candidateID).
		First(&app).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.ErrApplicationNotFound
	}
	return &app, err
}

func (r *applicationRepository) ListByCandidate(ctx context.Context, candidateID uuid.UUID) ([]model.Application, error) {
	var apps []model.Application
	err := r.db.WithContext(ctx).
		Preload("Job").
		Preload("Job.Stages").
		Preload("CurrentStage").
		Where("candidate_id = ?", candidateID).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) ListByJob(ctx context.Context, jobID uuid.UUID) ([]model.Application, error) {
	var apps []model.Application
	err := r.db.WithContext(ctx).
		Preload("Candidate").
		Preload("CurrentStage").
		Where("job_id = ?", jobID).
		Order("created_at DESC").
		Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) Update(ctx context.Context, app *model.Application) error {
	return r.db.WithContext(ctx).Save(app).Error
}

func (r *applicationRepository) RejectActiveApplications(ctx context.Context, jobID uuid.UUID) error {
	active := []model.ApplicationStatus{
		model.ApplicationStatusPending,
		model.ApplicationStatusInProgress,
	}
	return r.db.WithContext(ctx).
		Model(&model.Application{}).
		Where("job_id = ? AND status IN ?", jobID, active).
		Update("status", model.ApplicationStatusRejected).Error
}
