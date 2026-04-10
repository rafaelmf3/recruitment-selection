package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/api/handler"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/middleware"
	"recruitment-selection/internal/model"
	mockservice "recruitment-selection/internal/mock/service"
	"recruitment-selection/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupApplicationRouter(appSvc *mockservice.ApplicationService, userID interface{ String() string }, role model.UserRole) *gin.Engine {
	r := gin.New()
	h := handler.NewApplicationHandler(appSvc, 10)

	r.Use(middleware.InjectTestUser(testutil.CandidateID, role))

	r.POST("/api/v1/jobs/:id/apply", h.Apply)
	r.GET("/api/v1/applications", h.GetMyApplications)
	r.GET("/api/v1/jobs/:id/applications", h.GetJobApplications)
	r.PATCH("/api/v1/applications/:id/stage", h.AdvanceStage)
	r.PATCH("/api/v1/applications/:id/status", h.UpdateStatus)
	return r
}

// ---- GetMyApplications ------------------------------------------------------

func TestApplicationHandler_GetMyApplications_200(t *testing.T) {
	appSvc := new(mockservice.ApplicationService)
	router := setupApplicationRouter(appSvc, testutil.CandidateID, model.RoleCandidate)

	apps := []dto.ApplicationResponse{{Status: model.ApplicationStatusPending}}
	appSvc.On("GetMyApplications", anyCtx(), testutil.CandidateID).Return(apps, nil)

	w := doRequest(t, router, http.MethodGet, "/api/v1/applications", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Len(t, resp["data"], 1)
	appSvc.AssertExpectations(t)
}

// ---- AdvanceStage -----------------------------------------------------------

func TestApplicationHandler_AdvanceStage_200(t *testing.T) {
	appSvc := new(mockservice.ApplicationService)
	// Recruiter advances a candidate
	r := gin.New()
	r.Use(middleware.InjectTestUser(testutil.RecruiterID, model.RoleRecruiter))
	h := handler.NewApplicationHandler(appSvc, 10)
	r.PATCH("/api/v1/applications/:id/stage", h.AdvanceStage)

	stage := &dto.StageResponse{Name: "Technical Challenge", OrderIndex: 2}
	expected := &dto.ApplicationResponse{CurrentStage: stage}
	appSvc.On("AdvanceStage", anyCtx(), testutil.RecruiterID, testutil.ApplicationID).
		Return(expected, nil)

	w := doRequest(t, r, http.MethodPatch, "/api/v1/applications/"+testutil.ApplicationID.String()+"/stage", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	appSvc.AssertExpectations(t)
}

func TestApplicationHandler_AdvanceStage_422_LastStage(t *testing.T) {
	appSvc := new(mockservice.ApplicationService)
	r := gin.New()
	r.Use(middleware.InjectTestUser(testutil.RecruiterID, model.RoleRecruiter))
	h := handler.NewApplicationHandler(appSvc, 10)
	r.PATCH("/api/v1/applications/:id/stage", h.AdvanceStage)

	appSvc.On("AdvanceStage", anyCtx(), testutil.RecruiterID, testutil.ApplicationID).
		Return(nil, apierror.ErrNoNextStage)

	w := doRequest(t, r, http.MethodPatch, "/api/v1/applications/"+testutil.ApplicationID.String()+"/stage", nil)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	appSvc.AssertExpectations(t)
}

// ---- UpdateStatus -----------------------------------------------------------

func TestApplicationHandler_UpdateStatus_200_Accept(t *testing.T) {
	appSvc := new(mockservice.ApplicationService)
	r := gin.New()
	r.Use(middleware.InjectTestUser(testutil.RecruiterID, model.RoleRecruiter))
	h := handler.NewApplicationHandler(appSvc, 10)
	r.PATCH("/api/v1/applications/:id/status", h.UpdateStatus)

	body := dto.UpdateApplicationStatusRequest{Status: model.ApplicationStatusAccepted}
	expected := &dto.ApplicationResponse{Status: model.ApplicationStatusAccepted}

	appSvc.On("UpdateStatus", anyCtx(), testutil.RecruiterID, testutil.ApplicationID, model.ApplicationStatusAccepted).
		Return(expected, nil)

	w := doRequest(t, r, http.MethodPatch, "/api/v1/applications/"+testutil.ApplicationID.String()+"/status", body)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.ApplicationResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, model.ApplicationStatusAccepted, resp.Status)
	appSvc.AssertExpectations(t)
}

func TestApplicationHandler_UpdateStatus_403_NotOwner(t *testing.T) {
	appSvc := new(mockservice.ApplicationService)
	r := gin.New()
	r.Use(middleware.InjectTestUser(testutil.RecruiterID, model.RoleRecruiter))
	h := handler.NewApplicationHandler(appSvc, 10)
	r.PATCH("/api/v1/applications/:id/status", h.UpdateStatus)

	body := dto.UpdateApplicationStatusRequest{Status: model.ApplicationStatusRejected}
	appSvc.On("UpdateStatus", anyCtx(), testutil.RecruiterID, testutil.ApplicationID, model.ApplicationStatusRejected).
		Return(nil, apierror.ErrNotOwner)

	w := doRequest(t, r, http.MethodPatch, "/api/v1/applications/"+testutil.ApplicationID.String()+"/status", body)

	assert.Equal(t, http.StatusForbidden, w.Code)
	appSvc.AssertExpectations(t)
}
