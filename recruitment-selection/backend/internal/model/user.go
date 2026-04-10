package model

import (
	"time"

	"github.com/google/uuid"
)

// UserRole represents the role of a user in the system.
type UserRole string

const (
	RoleRecruiter UserRole = "recruiter"
	RoleCandidate UserRole = "candidate"
)

// User is the GORM model for the users table.
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name         string    `gorm:"not null"`
	Email        string    `gorm:"uniqueIndex;not null"`
	PasswordHash string    `gorm:"not null"`
	Role         UserRole  `gorm:"type:user_role;not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
