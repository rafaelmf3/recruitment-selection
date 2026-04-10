// Package apierror defines sentinel errors used across service and handler layers.
// Handlers map these errors to the appropriate HTTP status codes.
package apierror

import "errors"

var (
	// Auth errors
	ErrEmailAlreadyExists  = errors.New("email already in use")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")

	// Authorization errors
	ErrForbidden = errors.New("action not allowed for this role")
	ErrNotOwner  = errors.New("resource belongs to another user")

	// Job errors
	ErrJobNotFound        = errors.New("job not found")
	ErrJobNotAccepting    = errors.New("job is not accepting applications")
	ErrJobTerminal        = errors.New("job is in a terminal state and cannot be modified")
	ErrInvalidSalary      = errors.New("salary_min cannot be greater than salary_max")
	ErrInvalidTransition  = errors.New("invalid status transition")

	// Stage errors
	ErrStageNotFound    = errors.New("stage not found")
	ErrNoNextStage      = errors.New("candidate is already at the last stage")

	// Application errors
	ErrAlreadyApplied   = errors.New("already applied to this job")
	ErrApplicationNotFound = errors.New("application not found")
)
