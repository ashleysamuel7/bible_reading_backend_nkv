package integration

import (
	"os"
	"testing"

	"bible_reading_backend_nkv/database"
	"bible_reading_backend_nkv/server"

	"github.com/stretchr/testify/require"
)

// setupTestServer creates a test server instance
func setupTestServer(t *testing.T) server.Server {
	if os.Getenv("TEST_DB_DSN") == "" {
		t.Skip("TEST_DB_DSN not set - skipping integration tests")
		return nil
	}

	// Set TEST_DB_DSN as DB_DSN for database client
	originalDSN := os.Getenv("DB_DSN")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	defer os.Setenv("DB_DSN", originalDSN)

	db, err := database.NewDatabaseClient()
	require.NoError(t, err, "Failed to initialize test database")

	return server.NewEchoServer(db)
}

// checkDatabaseConnection verifies database is available
func checkDatabaseConnection(t *testing.T) bool {
	if os.Getenv("TEST_DB_DSN") == "" {
		return false
	}

	originalDSN := os.Getenv("DB_DSN")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	defer os.Setenv("DB_DSN", originalDSN)

	db, err := database.NewDatabaseClient()
	if err != nil {
		return false
	}

	return db.Ready()
}

// requireDatabase skips test if database is not available
func requireDatabase(t *testing.T) {
	if !checkDatabaseConnection(t) {
		t.Skip("Database not available - skipping integration test")
	}
}

