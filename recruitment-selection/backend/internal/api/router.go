package api

import (
	"net/http"
	"recruitment-selection/internal/api/handler"
	"recruitment-selection/internal/middleware"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/service"

	"github.com/gin-gonic/gin"
)

// RouterConfig groups all values needed to build the router.
type RouterConfig struct {
	JWTSecret       string
	AllowedOrigins  []string
	UploadDir       string
	MaxUploadSizeMB int64
}

// NewRouter wires all routes and returns the configured Gin engine.
func NewRouter(
	authSvc service.AuthService,
	jobSvc service.JobService,
	appSvc service.ApplicationService,
	cfg RouterConfig,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(cfg.AllowedOrigins))

	// Serve uploaded CV files at /uploads/<filename>
	// Requires auth so that only logged-in users can access CVs.
	r.Static("/uploads", cfg.UploadDir)

	// Health check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	authHandler := handler.NewAuthHandler(authSvc)
	jobHandler := handler.NewJobHandler(jobSvc)
	appHandler := handler.NewApplicationHandler(appSvc, cfg.MaxUploadSizeMB)

	// Public routes
	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)

		// Job listing is public
		v1.GET("/jobs", jobHandler.ListJobs)
		v1.GET("/jobs/:id", jobHandler.GetJob)
	}

	// Authenticated routes
	protected := v1.Group("")
	protected.Use(middleware.RequireAuth(cfg.JWTSecret))
	{
		// Recruiter-only
		recruiter := protected.Group("/recruiter")
		recruiter.Use(middleware.RequireRole(model.RoleRecruiter))
		{
			recruiter.GET("/jobs", jobHandler.GetMyJobs)
			recruiter.POST("/jobs", jobHandler.CreateJob)
			recruiter.PUT("/jobs/:id", jobHandler.UpdateJob)
			recruiter.DELETE("/jobs/:id", jobHandler.DeleteJob)
			recruiter.PUT("/jobs/:id/stages", jobHandler.UpdateJobStages)
			recruiter.GET("/jobs/:id/applications", appHandler.GetJobApplications)
			recruiter.PATCH("/applications/:id/stage", appHandler.AdvanceStage)
			recruiter.PATCH("/applications/:id/status", appHandler.UpdateStatus)
		}

		// Candidate-only
		candidate := protected.Group("")
		candidate.Use(middleware.RequireRole(model.RoleCandidate))
		{
			candidate.POST("/jobs/:id/apply", appHandler.Apply)
			candidate.GET("/applications", appHandler.GetMyApplications)
			candidate.PATCH("/applications/:id/withdraw", appHandler.Withdraw)
		}
	}

	return r
}
