package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/api/handler"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/model"
	mockservice "recruitment-selection/internal/mock/service"
	"recruitment-selection/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

func setupAuthRouter(authSvc *mockservice.AuthService) *gin.Engine {
	r := gin.New()
	h := handler.NewAuthHandler(authSvc)
	r.POST("/api/v1/auth/register", h.Register)
	r.POST("/api/v1/auth/login", h.Login)
	return r
}

// ---- Register ---------------------------------------------------------------

func TestAuthHandler_Register_201(t *testing.T) {
	authSvc := new(mockservice.AuthService)
	router := setupAuthRouter(authSvc)

	body := dto.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@company.com",
		Password: "secret123",
		Role:     model.RoleRecruiter,
	}
	expected := &dto.LoginResponse{
		Token: "jwt.token.here",
		User: dto.UserResponse{
			ID:    testutil.RecruiterID.String(),
			Name:  body.Name,
			Email: body.Email,
			Role:  body.Role,
		},
	}

	authSvc.On("Register", anyCtx(), body).Return(expected, nil)

	w := doRequest(t, router, http.MethodPost, "/api/v1/auth/register", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp dto.LoginResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, expected.User.Email, resp.User.Email)
	assert.NotEmpty(t, resp.Token)
	authSvc.AssertExpectations(t)
}

func TestAuthHandler_Register_400_InvalidBody(t *testing.T) {
	authSvc := new(mockservice.AuthService)
	router := setupAuthRouter(authSvc)

	// Missing required fields
	w := doRequest(t, router, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"email": "not-an-email",
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	authSvc.AssertNotCalled(t, "Register")
}

func TestAuthHandler_Register_409_EmailConflict(t *testing.T) {
	authSvc := new(mockservice.AuthService)
	router := setupAuthRouter(authSvc)

	body := dto.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@company.com",
		Password: "secret123",
		Role:     model.RoleRecruiter,
	}

	authSvc.On("Register", anyCtx(), body).Return((*dto.LoginResponse)(nil), apierror.ErrEmailAlreadyExists)

	w := doRequest(t, router, http.MethodPost, "/api/v1/auth/register", body)

	assert.Equal(t, http.StatusConflict, w.Code)
	authSvc.AssertExpectations(t)
}

// ---- Login ------------------------------------------------------------------

func TestAuthHandler_Login_200(t *testing.T) {
	authSvc := new(mockservice.AuthService)
	router := setupAuthRouter(authSvc)

	body := dto.LoginRequest{Email: "alice@company.com", Password: "secret123"}
	expected := &dto.LoginResponse{
		Token: "jwt.token.here",
		User: dto.UserResponse{
			ID:    testutil.RecruiterID.String(),
			Email: body.Email,
			Role:  model.RoleRecruiter,
		},
	}

	authSvc.On("Login", anyCtx(), body).Return(expected, nil)

	w := doRequest(t, router, http.MethodPost, "/api/v1/auth/login", body)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.LoginResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotEmpty(t, resp.Token)
	authSvc.AssertExpectations(t)
}

func TestAuthHandler_Login_400_InvalidBody(t *testing.T) {
	authSvc := new(mockservice.AuthService)
	router := setupAuthRouter(authSvc)

	w := doRequest(t, router, http.MethodPost, "/api/v1/auth/login", map[string]string{})

	assert.Equal(t, http.StatusBadRequest, w.Code)
	authSvc.AssertNotCalled(t, "Login")
}

func TestAuthHandler_Login_401_BadCredentials(t *testing.T) {
	authSvc := new(mockservice.AuthService)
	router := setupAuthRouter(authSvc)

	body := dto.LoginRequest{Email: "alice@company.com", Password: "wrong"}
	authSvc.On("Login", anyCtx(), body).Return(nil, apierror.ErrInvalidCredentials)

	w := doRequest(t, router, http.MethodPost, "/api/v1/auth/login", body)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	authSvc.AssertExpectations(t)
}
