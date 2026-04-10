package repository

import (
	"context"
	"errors"
	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// jobRepository is the PostgreSQL implementation of JobRepository.
type jobRepository struct {
	db *gorm.DB
}

// NewJobRepository returns a new JobRepository backed by the given DB connection.
func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *model.Job) error {
	return r.db.WithContext(ctx).Create(job).Error
}

func (r *jobRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Job, error) {
	var job model.Job
	err := r.db.WithContext(ctx).
		Preload("Recruiter").
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		First(&job, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.ErrJobNotFound
	}
	return &job, err
}

func (r *jobRepository) FindByRecruiter(ctx context.Context, recruiterID uuid.UUID, filter dto.RecruiterJobFilter) ([]model.Job, error) {
	var jobs []model.Job

	query := r.db.WithContext(ctx).Where("recruiter_id = ?", recruiterID)

	if filter.Q != "" {
		like := "%" + filter.Q + "%"
		query = query.Where("title ILIKE ? OR company ILIKE ?", like, like)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	err := query.
		Preload("Recruiter").
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) List(ctx context.Context, filter dto.JobFilter) ([]model.Job, int64, error) {
	var jobs []model.Job
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Job{}).Where("status = ?", model.JobStatusOpen)

	if filter.Q != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Q+"%")
	}
	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}
	if filter.SalaryMin != nil {
		query = query.Where("salary_min >= ?", *filter.SalaryMin)
	}
	if filter.SalaryMax != nil {
		query = query.Where("salary_max <= ?", *filter.SalaryMax)
	}

	// Count before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filter.Page - 1) * filter.Limit
	err := query.
		Preload("Recruiter").
		Preload("Stages", func(db *gorm.DB) *gorm.DB {
			return db.Order("order_index ASC")
		}).
		Order("created_at DESC").
		Offset(offset).
		Limit(filter.Limit).
		Find(&jobs).Error

	return jobs, total, err
}

func (r *jobRepository) Update(ctx context.Context, job *model.Job) error {
	return r.db.WithContext(ctx).Save(job).Error
}

func (r *jobRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Job{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apierror.ErrJobNotFound
	}
	return nil
}

func (r *jobRepository) ReplaceStages(ctx context.Context, jobID uuid.UUID, stages []model.JobStage) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete all existing stages for this job
		if err := tx.Delete(&model.JobStage{}, "job_id = ?", jobID).Error; err != nil {
			return err
		}
		if len(stages) == 0 {
			return nil
		}
		return tx.Create(&stages).Error
	})
}

func (r *jobRepository) FindStageByID(ctx context.Context, id uuid.UUID) (*model.JobStage, error) {
	var stage model.JobStage
	err := r.db.WithContext(ctx).First(&stage, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.ErrStageNotFound
	}
	return &stage, err
}

func (r *jobRepository) FindStagesByJobID(ctx context.Context, jobID uuid.UUID) ([]model.JobStage, error) {
	var stages []model.JobStage
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("order_index ASC").
		Find(&stages).Error
	return stages, err
}
