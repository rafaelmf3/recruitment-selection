//go:build integration

package testutil

import (
	"fmt"
	"os"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NewTestDB opens a connection to the test database defined in docker-compose
// (db_test service on port 5433). The caller must ensure the container is
// running before executing integration tests:
//
//	docker compose up -d db_test
//	go test ./internal/... -tags=integration -v
func NewTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := testDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("testutil: failed to connect to test database: %v", err)
	}

	return db
}

// CleanTables truncates all application tables in dependency-safe order.
// Call this in t.Cleanup() to leave the DB clean between tests.
func CleanTables(t *testing.T, db *gorm.DB) {
	t.Helper()
	tables := []string{"applications", "job_stages", "jobs", "users"}
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			t.Fatalf("testutil: failed to truncate %s: %v", table, err)
		}
	}
}

func testDSN() string {
	host := getenv("TEST_DB_HOST", "localhost")
	port := getenv("TEST_DB_PORT", "5433")
	user := getenv("TEST_DB_USER", "postgres")
	pass := getenv("TEST_DB_PASSWORD", "postgres")
	name := getenv("TEST_DB_NAME", "recruitment_selection_test")
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
