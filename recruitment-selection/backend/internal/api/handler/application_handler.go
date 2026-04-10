package handler

import (
	"errors"
	"fmt"
	"net/http"
	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/dto"
	"recruitment-selection/internal/middleware"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ApplicationHandler handles application endpoints.
type ApplicationHandler struct {
	appSvc          service.ApplicationService
	maxUploadSizeMB int64
}

// NewApplicationHandler returns a new ApplicationHandler.
// maxUploadSizeMB is the maximum allowed CV file size in megabytes.
func NewApplicationHandler(appSvc service.ApplicationService, maxUploadSizeMB int64) *ApplicationHandler {
	return &ApplicationHandler{appSvc: appSvc, maxUploadSizeMB: maxUploadSizeMB}
}

// Apply godoc
// POST /api/v1/jobs/:id/apply
func (h *ApplicationHandler) Apply(c *gin.Context) {
	candidateID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	// Enforce max upload size before reading the form
	maxBytes := h.maxUploadSizeMB << 20 // MB to bytes
	if err := c.Request.ParseMultipartForm(maxBytes); err != nil {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{
			"error": fmt.Sprintf("file too large: max %d MB allowed", h.maxUploadSizeMB),
		})
		return
	}

	coverLetter := c.PostForm("cover_letter")
	cvFile, err := c.FormFile("cv")
	if err != nil || cvFile == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "curriculum (cv) é obrigatório"})
		return
	}

	resp, err := h.appSvc.Apply(c.Request.Context(), candidateID, jobID, coverLetter, cvFile)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrJobNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrJobNotAccepting):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrAlreadyApplied):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// GetMyApplications godoc
// GET /api/v1/applications
func (h *ApplicationHandler) GetMyApplications(c *gin.Context) {
	candidateID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	apps, err := h.appSvc.GetMyApplications(c.Request.Context(), candidateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if apps == nil {
		apps = []dto.ApplicationResponse{}
	}
	c.JSON(http.StatusOK, apps)
}

// GetJobApplications godoc
// GET /api/v1/jobs/:id/applications
func (h *ApplicationHandler) GetJobApplications(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
		return
	}

	apps, err := h.appSvc.GetJobApplications(c.Request.Context(), recruiterID, jobID)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	if apps == nil {
		apps = []dto.ApplicationResponse{}
	}
	c.JSON(http.StatusOK, apps)
}

// Withdraw godoc
// PATCH /api/v1/applications/:id/withdraw
func (h *ApplicationHandler) Withdraw(c *gin.Context) {
	candidateID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	resp, err := h.appSvc.Withdraw(c.Request.Context(), candidateID, appID)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrApplicationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// AdvanceStage godoc
// PATCH /api/v1/applications/:id/stage
func (h *ApplicationHandler) AdvanceStage(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	resp, err := h.appSvc.AdvanceStage(c.Request.Context(), recruiterID, appID)
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrApplicationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNoNextStage):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateStatus godoc
// PATCH /api/v1/applications/:id/status
func (h *ApplicationHandler) UpdateStatus(c *gin.Context) {
	recruiterID := c.MustGet(middleware.ContextKeyUserID).(uuid.UUID)

	appID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	var req dto.UpdateApplicationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.appSvc.UpdateStatus(c.Request.Context(), recruiterID, appID, model.ApplicationStatus(req.Status))
	if err != nil {
		switch {
		case errors.Is(err, apierror.ErrApplicationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, apierror.ErrNotOwner):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}
