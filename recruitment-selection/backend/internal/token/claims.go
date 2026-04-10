// Package token defines the JWT claims struct used by both the auth service
// (token generation) and the auth middleware (token validation).
// Keeping it here avoids circular imports between those two packages.
package token

import (
	"recruitment-selection/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the payload encoded inside every JWT issued by the system.
type Claims struct {
	UserID string         `json:"user_id"`
	Email  string         `json:"email"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}
