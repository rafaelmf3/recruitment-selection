//go:build integration

package repository_test

import (
	"context"
	"testing"

	"recruitment-selection/internal/apierror"
	"recruitment-selection/internal/model"
	"recruitment-selection/internal/repository"
	"recruitment-selection/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create_And_FindByEmail(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	repo := repository.NewUserRepository(db)

	user := &model.User{
		ID:           uuid.New(),
		Name:         "Alice",
		Email:        "alice@company.com",
		PasswordHash: "hashed",
		Role:         model.RoleRecruiter,
	}

	err := repo.Create(context.Background(), user)
	require.NoError(t, err)

	found, err := repo.FindByEmail(context.Background(), user.Email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
	assert.Equal(t, user.Email, found.Email)
	assert.Equal(t, user.Role, found.Role)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	repo := repository.NewUserRepository(db)

	_, err := repo.FindByEmail(context.Background(), "nobody@example.com")
	assert.ErrorIs(t, err, apierror.ErrUserNotFound)
}

func TestUserRepository_Create_DuplicateEmail(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	repo := repository.NewUserRepository(db)

	user := &model.User{
		ID:           uuid.New(),
		Name:         "Alice",
		Email:        "alice@company.com",
		PasswordHash: "hashed",
		Role:         model.RoleRecruiter,
	}
	require.NoError(t, repo.Create(context.Background(), user))

	duplicate := &model.User{
		ID:           uuid.New(),
		Name:         "Alice Clone",
		Email:        user.Email, // same email
		PasswordHash: "hashed2",
		Role:         model.RoleCandidate,
	}
	err := repo.Create(context.Background(), duplicate)
	assert.Error(t, err, "duplicate email must return an error")
}

func TestUserRepository_FindByID(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	repo := repository.NewUserRepository(db)

	user := &model.User{
		ID:           uuid.New(),
		Name:         "Bob",
		Email:        "bob@gmail.com",
		PasswordHash: "hashed",
		Role:         model.RoleCandidate,
	}
	require.NoError(t, repo.Create(context.Background(), user))

	found, err := repo.FindByID(context.Background(), user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, found.ID)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := testutil.NewTestDB(t)
	t.Cleanup(func() { testutil.CleanTables(t, db) })

	repo := repository.NewUserRepository(db)

	_, err := repo.FindByID(context.Background(), uuid.New())
	assert.ErrorIs(t, err, apierror.ErrUserNotFound)
}
