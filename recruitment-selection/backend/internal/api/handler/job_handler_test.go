package handler_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"recruitment-selection/internal/api/handler"
	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/middleware"
	"recruitment-selection/internal/model"
	mockservice "recruitment-selection/internal/mock/service"
	"recruitment-selection/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupJobRouter(jobSvc *mockservice.JobService, role model.UserRole) *gin.Engine {
	r := gin.New()
	h := handler.NewJobHandler(jobSvc)

	r.GET("/api/v1/jobs", h.ListJobs)
	r.GET("/api/v1/jobs/:id", h.GetJob)

	auth := r.Group("/api/v1/jobs")
	auth.Use(middleware.InjectTestUser(testutil.RecruiterID, role))
	{
		auth.POST("", h.CreateJob)
		auth.PUT("/:id", h.UpdateJob)
		auth.DELETE("/:id", h.DeleteJob)
		auth.PUT("/:id/stages", h.UpdateJobStages)
	}
	return r
}

// ---- ListJobs ---------------------------------------------------------------

func TestJobHandler_ListJobs_200(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	jobs := []dto.JobResponse{{Title: "Backend Engineer"}}
	jobSvc.On("ListJobs", anyCtx(), dto.JobFilter{Page: 1, Limit: 20}).
		Return(jobs, int64(1), nil)

	w := doRequest(t, router, http.MethodGet, "/api/v1/jobs?page=1&limit=20", nil)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp dto.PaginatedResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(1), resp.Total)
	assert.Equal(t, 1, resp.Page)
	assert.Equal(t, 20, resp.Limit)
	jobSvc.AssertExpectations(t)
}

// ---- GetJob -----------------------------------------------------------------

func TestJobHandler_GetJob_200(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	job := &dto.JobResponse{Title: "Backend Engineer"}
	jobSvc.On("GetJobByID", anyCtx(), testutil.JobID).Return(job, nil)

	w := doRequest(t, router, http.MethodGet, "/api/v1/jobs/"+testutil.JobID.String(), nil)

	assert.Equal(t, http.StatusOK, w.Code)
	jobSvc.AssertExpectations(t)
}

func TestJobHandler_GetJob_404(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	jobSvc.On("GetJobByID", anyCtx(), testutil.JobID).Return(nil, apierror.ErrJobNotFound)

	w := doRequest(t, router, http.MethodGet, "/api/v1/jobs/"+testutil.JobID.String(), nil)

	assert.Equal(t, http.StatusNotFound, w.Code)
	jobSvc.AssertExpectations(t)
}

func TestJobHandler_GetJob_400_InvalidUUID(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	w := doRequest(t, router, http.MethodGet, "/api/v1/jobs/not-a-uuid", nil)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	jobSvc.AssertNotCalled(t, "GetJobByID")
}

// ---- CreateJob --------------------------------------------------------------

func TestJobHandler_CreateJob_201(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	body := dto.CreateJobRequest{
		Title:       "Backend Engineer",
		Description: "Build APIs.",
		SalaryMin:   testutil.SalaryMin,
		SalaryMax:   testutil.SalaryMax,
	}
	expected := &dto.JobResponse{Title: body.Title}
	jobSvc.On("CreateJob", anyCtx(), testutil.RecruiterID, body).Return(expected, nil)

	w := doRequest(t, router, http.MethodPost, "/api/v1/jobs", body)

	assert.Equal(t, http.StatusCreated, w.Code)
	jobSvc.AssertExpectations(t)
}

func TestJobHandler_CreateJob_400_InvalidSalary(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	min := 9000.0
	max := 1000.0
	body := dto.CreateJobRequest{
		Title:       "Backend Engineer",
		Description: "Build APIs.",
		SalaryMin:   &min,
		SalaryMax:   &max,
	}
	jobSvc.On("CreateJob", anyCtx(), testutil.RecruiterID, body).Return(nil, apierror.ErrInvalidSalary)

	w := doRequest(t, router, http.MethodPost, "/api/v1/jobs", body)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	jobSvc.AssertExpectations(t)
}

// ---- DeleteJob --------------------------------------------------------------

func TestJobHandler_DeleteJob_204(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	jobSvc.On("DeleteJob", anyCtx(), testutil.RecruiterID, testutil.JobID).Return(nil)

	w := doRequest(t, router, http.MethodDelete, "/api/v1/jobs/"+testutil.JobID.String(), nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	jobSvc.AssertExpectations(t)
}

func TestJobHandler_DeleteJob_403_NotOwner(t *testing.T) {
	jobSvc := new(mockservice.JobService)
	router := setupJobRouter(jobSvc, model.RoleRecruiter)

	jobSvc.On("DeleteJob", anyCtx(), testutil.RecruiterID, testutil.JobID).Return(apierror.ErrNotOwner)

	w := doRequest(t, router, http.MethodDelete, "/api/v1/jobs/"+testutil.JobID.String(), nil)

	assert.Equal(t, http.StatusForbidden, w.Code)
	jobSvc.AssertExpectations(t)
}
