package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/repository"
	"recruitment-selection/internal/token"
)

// authService is the concrete implementation of AuthService.
type authService struct {
	userRepo           repository.UserRepository
	jwtSecret          string
	jwtExpirationHours int
}

// NewAuthService returns a new AuthService.
func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: 24,
	}
}

// NewAuthServiceWithExpiration returns a new AuthService with a custom token TTL.
// Useful in tests that need shorter-lived tokens.
func NewAuthServiceWithExpiration(userRepo repository.UserRepository, jwtSecret string, expirationHours int) AuthService {
	return &authService{
		userRepo:           userRepo,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: expirationHours,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Reject if email is already taken
	_, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, apierror.ErrEmailAlreadyExists
	}
	if !errors.Is(err, apierror.ErrUserNotFound) {
		return nil, fmt.Errorf("register: check email: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("register: hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("register: create user: %w", err)
	}

	t, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("register: generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: t,
		User:  *toUserResponse(user),
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		// Return a generic error to prevent user enumeration
		return nil, apierror.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apierror.ErrInvalidCredentials
	}

	t, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("login: generate token: %w", err)
	}

	return &dto.LoginResponse{
		Token: t,
		User:  *toUserResponse(user),
	}, nil
}

// generateToken signs a JWT with the user's ID, email and role.
func (s *authService) generateToken(user *model.User) (string, error) {
	claims := token.Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.jwtSecret))
}

// toUserResponse maps a model.User to the public DTO (no password hash).
func toUserResponse(u *model.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:    u.ID.String(),
		Name:  u.Name,
		Email: u.Email,
		Role:  u.Role,
	}
}
