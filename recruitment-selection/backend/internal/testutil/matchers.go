package testutil

import (
	"recruitment-selection/internal/model"

	"github.com/stretchr/testify/mock"
)

// AnyUser returns a testify matcher that accepts any *model.User argument.
// Used when the exact user struct is constructed inside the service and
// cannot be predicted precisely (e.g., generated UUID, bcrypt hash).
func AnyUser() interface{} {
	return mock.MatchedBy(func(u *model.User) bool {
		return u != nil
	})
}

// AnyApplication returns a testify matcher that accepts any *model.Application.
func AnyApplication() interface{} {
	return mock.MatchedBy(func(a *model.Application) bool {
		return a != nil
	})
}

// AnyJob returns a testify matcher that accepts any *model.Job.
func AnyJob() interface{} {
	return mock.MatchedBy(func(j *model.Job) bool {
		return j != nil
	})
}
