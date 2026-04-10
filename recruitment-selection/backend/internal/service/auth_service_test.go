package service_test

import (
	"context"
	"errors"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	mockrepo "recruitment-selection/internal/mock/repository"
	"recruitment-selection/internal/service"
	"recruitment-selection/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

const testJWTSecret = "test-secret-key"

func newAuthService(userRepo *mockrepo.UserRepository) service.AuthService {
	return service.NewAuthService(userRepo, testJWTSecret)
}

// ---- Register ---------------------------------------------------------------

func TestAuthService_Register_Success_Recruiter(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	req := dto.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@company.com",
		Password: "secret123",
		Role:     model.RoleRecruiter,
	}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(nil, apierror.ErrUserNotFound)
	userRepo.On("Create", context.Background(), testutil.AnyUser()).
		Return(nil)

	resp, err := svc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token, "JWT token must be present after register")
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, model.RoleRecruiter, resp.User.Role)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_Success_Candidate(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	req := dto.RegisterRequest{
		Name:     "Bob",
		Email:    "bob@gmail.com",
		Password: "secret123",
		Role:     model.RoleCandidate,
	}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(nil, apierror.ErrUserNotFound)
	userRepo.On("Create", context.Background(), testutil.AnyUser()).
		Return(nil)

	resp, err := svc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, model.RoleCandidate, resp.User.Role)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_EmailAlreadyExists(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	existingUser := testutil.NewRecruiter()
	req := dto.RegisterRequest{
		Name:     "Alice Duplicate",
		Email:    existingUser.Email,
		Password: "secret123",
		Role:     model.RoleRecruiter,
	}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(existingUser, nil)

	resp, err := svc.Register(context.Background(), req)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrEmailAlreadyExists))
	userRepo.AssertExpectations(t)
}

func TestAuthService_Register_PasswordIsHashed(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	plainPassword := "my-plain-password"
	req := dto.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@company.com",
		Password: plainPassword,
		Role:     model.RoleRecruiter,
	}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(nil, apierror.ErrUserNotFound)

	var capturedUser *model.User
	userRepo.On("Create", context.Background(), testutil.AnyUser()).
		Run(func(args mock.Arguments) {
			capturedUser = args.Get(1).(*model.User)
		}).
		Return(nil)

	_, err := svc.Register(context.Background(), req)

	require.NoError(t, err)
	assert.NotEqual(t, plainPassword, capturedUser.PasswordHash,
		"password must not be stored in plain text")
	assert.NotEmpty(t, capturedUser.PasswordHash)
}

// ---- Login ------------------------------------------------------------------

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	user := testutil.NewRecruiter()
	// Use a real bcrypt hash of "secret123" so the service can verify it.
	user.PasswordHash = mustBcrypt("secret123")

	req := dto.LoginRequest{Email: user.Email, Password: "secret123"}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(user, nil)

	resp, err := svc.Login(context.Background(), req)

	require.NoError(t, err)
	assert.NotEmpty(t, resp.Token, "JWT token must be present")
	assert.Equal(t, user.Email, resp.User.Email)
	assert.Equal(t, user.Role, resp.User.Role)
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	req := dto.LoginRequest{Email: "nobody@example.com", Password: "secret123"}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(nil, apierror.ErrUserNotFound)

	resp, err := svc.Login(context.Background(), req)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrInvalidCredentials))
	userRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	userRepo := new(mockrepo.UserRepository)
	svc := newAuthService(userRepo)

	user := testutil.NewRecruiter()
	user.PasswordHash = mustBcrypt("correct-password")

	req := dto.LoginRequest{Email: user.Email, Password: "wrong-password"}

	userRepo.On("FindByEmail", context.Background(), req.Email).
		Return(user, nil)

	resp, err := svc.Login(context.Background(), req)

	assert.Nil(t, resp)
	assert.True(t, errors.Is(err, apierror.ErrInvalidCredentials))
	userRepo.AssertExpectations(t)
}

// ---- Helpers ----------------------------------------------------------------

// mustBcrypt generates a bcrypt hash for use in test fixtures.
func mustBcrypt(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
