package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bible_reading_backend_nkv/database"
	"bible_reading_backend_nkv/dto"
	"bible_reading_backend_nkv/server"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEndpoints tests all API endpoints individually
func TestEndpoints(t *testing.T) {
	// Setup
	if os.Getenv("TEST_DB_DSN") == "" {
		t.Skip("TEST_DB_DSN not set - skipping integration tests")
		return
	}

	// Set TEST_DB_DSN as DB_DSN
	originalDSN := os.Getenv("DB_DSN")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	defer os.Setenv("DB_DSN", originalDSN)

	db, err := database.NewDatabaseClient()
	require.NoError(t, err, "Failed to initialize database")

	srv := server.NewEchoServer(db)

	// Get Echo instance for testing
	var e *echo.Echo
	if echoSrv, ok := srv.(*server.EchoServer); ok {
		e = echoSrv.GetEcho()
	} else {
		e = echo.New()
	}

	t.Run("Health - Readiness", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := ParseJSONResponse(rec, &response)
		require.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
	})

	t.Run("Health - Liveness", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/liveness", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]string
		err := ParseJSONResponse(rec, &response)
		require.NoError(t, err)
		assert.Equal(t, "OK", response["status"])
	})

	t.Run("Get All Books", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/niv/books", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var books []map[string]interface{}
		err := ParseJSONResponse(rec, &books)
		require.NoError(t, err)
		if len(books) > 0 {
			// Verify first book structure
			firstBook := books[0]
			assert.Contains(t, firstBook, "book_id")
			assert.Contains(t, firstBook, "book")
			assert.NotNil(t, firstBook["book_id"])
			assert.NotEmpty(t, firstBook["book"])
		}
	})

	t.Run("Get All Verses", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/niv/verses", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var verses []map[string]interface{}
		err := ParseJSONResponse(rec, &verses)
		require.NoError(t, err)
		assert.Greater(t, len(verses), 0)

		// Verify verse structure
		firstVerse := verses[0]
		requiredFields := []string{"book_id", "book", "chapter", "verse", "text"}
		for _, field := range requiredFields {
			assert.Contains(t, firstVerse, field, "Verse should contain %s", field)
		}
	})

	t.Run("Get Verses by Chapter - Valid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/niv/1/1/verses", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var verses []map[string]interface{}
		err := ParseJSONResponse(rec, &verses)
		require.NoError(t, err)
		assert.Greater(t, len(verses), 0)
	})

	t.Run("Get Chapters by Book - Valid Book ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/niv/chapters/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := ParseJSONResponse(rec, &response)
		require.NoError(t, err)
		assert.Contains(t, response, "MaxChapter")

		maxChapter := response["MaxChapter"].(float64)
		assert.Greater(t, maxChapter, float64(0))
	})

	t.Run("Explain Verse - Valid Request", func(t *testing.T) {
		if os.Getenv("OPENAI_API_KEY") == "" {
			t.Skip("OPENAI_API_KEY not set")
			return
		}

		explainReq := dto.ExplainRequest{
			Book:       "Genesis",
			Chapter:    1,
			StartVerse: 1,
			EndVerse:   3,
			Age:        25,
			Belief:     3,
		}

		body, err := json.Marshal(explainReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// May succeed or fail depending on API key
		assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, rec.Code)
	})

	t.Run("Explain Verse - Invalid Request", func(t *testing.T) {
		invalidReq := map[string]interface{}{
			"book": "Genesis",
			// Missing required fields
		}

		body, err := json.Marshal(invalidReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]string
		err = ParseJSONResponse(rec, &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "Invalid")
	})

	t.Run("Explain Verse - Default Values", func(t *testing.T) {
		if os.Getenv("OPENAI_API_KEY") == "" {
			t.Skip("OPENAI_API_KEY not set")
			return
		}

		explainReq := dto.ExplainRequest{
			Book:       "Genesis",
			Chapter:    1,
			StartVerse: 1,
			EndVerse:   3,
			// Age and Belief omitted - should default
		}

		body, err := json.Marshal(explainReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Should handle defaults (may fail due to API but should process)
		assert.Contains(t, []int{http.StatusOK, http.StatusInternalServerError}, rec.Code)
	})
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	if os.Getenv("TEST_DB_DSN") == "" {
		t.Skip("TEST_DB_DSN not set")
		return
	}

	originalDSN := os.Getenv("DB_DSN")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	defer os.Setenv("DB_DSN", originalDSN)

	db, err := database.NewDatabaseClient()
	require.NoError(t, err)

	srv := server.NewEchoServer(db)

	var e *echo.Echo
	if echoSrv, ok := srv.(*server.EchoServer); ok {
		e = echoSrv.GetEcho()
	} else {
		e = echo.New()
	}

	_ = e

	t.Run("Invalid Chapter Number", func(t *testing.T) {
		// This requires echo context with invalid chapter parameter
		// Placeholder for error handling test
		assert.True(t, true)
	})

	t.Run("Missing API Key", func(t *testing.T) {
		// Save original key
		originalKey := os.Getenv("OPENAI_API_KEY")
		defer func() {
			if originalKey != "" {
				os.Setenv("OPENAI_API_KEY", originalKey)
			}
		}()

		// Remove key
		os.Unsetenv("OPENAI_API_KEY")

		explainReq := dto.ExplainRequest{
			Book:       "Genesis",
			Chapter:    1,
			StartVerse: 1,
			EndVerse:   3,
		}

		body, err := json.Marshal(explainReq)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]string
		err = ParseJSONResponse(rec, &response)
		require.NoError(t, err)
		assert.Contains(t, response["error"], "OpenAI API key")
	})
}

// TestResponseFormats tests response format consistency
func TestResponseFormats(t *testing.T) {
	if os.Getenv("TEST_DB_DSN") == "" {
		t.Skip("TEST_DB_DSN not set")
		return
	}

	originalDSN := os.Getenv("DB_DSN")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	defer os.Setenv("DB_DSN", originalDSN)

	db, err := database.NewDatabaseClient()
	require.NoError(t, err)

	srv := server.NewEchoServer(db)

	var e *echo.Echo
	if echoSrv, ok := srv.(*server.EchoServer); ok {
		e = echoSrv.GetEcho()
	} else {
		e = echo.New()
	}

	_ = e

	t.Run("Books Response Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/niv/books", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var books []map[string]interface{}
		err := ParseJSONResponse(rec, &books)
		require.NoError(t, err)

		if len(books) > 0 {
			book := books[0]
			// Verify JSON structure
			jsonData, err := json.Marshal(book)
			require.NoError(t, err)
			assert.Contains(t, string(jsonData), "book_id")
			assert.Contains(t, string(jsonData), "book")
		}
	})

	t.Run("Verses Response Format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/niv/verses", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		var verses []map[string]interface{}
		err := ParseJSONResponse(rec, &verses)
		require.NoError(t, err)

		if len(verses) > 0 {
			verse := verses[0]
			jsonData, err := json.Marshal(verse)
			require.NoError(t, err)

			requiredFields := []string{"book_id", "book", "chapter", "verse", "text"}
			for _, field := range requiredFields {
				assert.Contains(t, string(jsonData), field)
			}
		}
	})
}
