package handler

import (
	"errors"
	"net/http"
	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/middleware"
	"recruitment-selection/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// JobHandler handles job posting endpoints.
type JobHandler struct {
	jobSvc service.JobService
}

// NewJobHandler returns a new JobHandler.
func NewJobHandler(jobSvc service.JobService) *JobHandler {
	return &JobHandler{jobSvc: jobSvc}
}

// GetMyJobs godoc
// GET /api/v1/recruiter/jobs
func (h *JobHandler) GetMyJobs(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var filter dto.RecruiterJobFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, err := h.jobSvc.GetMyJobs(c.Request.Context(), recruiterID, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

// ListJobs godoc
// GET /api/v1/jobs
func (h *JobHandler) ListJobs(c *gin.Context) {
	var filter dto.JobFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jobs, total, err := h.jobSvc.ListJobs(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, dto.NewPaginated(jobs, total, filter.Page, filter.Limit))
}

// GetJob godoc
// GET /api/v1/jobs/:id
func (h *JobHandler) GetJob(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	job, err := h.jobSvc.GetJobByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, apierror.ErrJobNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, job)
}

// CreateJob godoc
// POST /api/v1/jobs
func (h *JobHandler) CreateJob(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	var req dto.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobSvc.CreateJob(c.Request.Context(), recruiterID, req)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrInvalidSalary):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, job)
}

// UpdateJob godoc
// PUT /api/v1/jobs/:id
func (h *JobHandler) UpdateJob(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	var req dto.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobSvc.UpdateJob(c.Request.Context(), recruiterID, id, req)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrJobNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrJobTerminal):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrInvalidTransition):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, job)
}

// DeleteJob godoc
// DELETE /api/v1/jobs/:id
func (h *JobHandler) DeleteJob(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	if err := h.jobSvc.DeleteJob(c.Request.Context(), recruiterID, id); err != nil {
		switch {
		case errors.Is(err, apierror.ErrJobNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateJobStages godoc
// PUT /api/v1/jobs/:id/stages
func (h *JobHandler) UpdateJobStages(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	var req dto.UpdateJobStagesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stages, err := h.jobSvc.UpdateJobStages(c.Request.Context(), recruiterID, id, req)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrJobNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, stages)
}
