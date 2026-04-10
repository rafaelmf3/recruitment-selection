package repository

import (
	"context"
	"errors"
	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userRepository is the PostgreSQL implementation of UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a new UserRepository backed by the given DB connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.ErrUserNotFound
	}
	return &user, err
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apierror.ErrUserNotFound
	}
	return &user, err
}
