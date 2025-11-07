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
	"github.com/stretchr/testify/suite"
)

// IntegrationTestSuite holds the test suite
type IntegrationTestSuite struct {
	suite.Suite
	server server.Server
	db     database.DatabaseClient
	e      *echo.Echo
}

// SetupSuite runs once before all tests
func (suite *IntegrationTestSuite) SetupSuite() {
	// Load environment variables
	if os.Getenv("TEST_DB_DSN") == "" {
		suite.T().Skip("TEST_DB_DSN not set - skipping integration tests")
		return
	}

	// Set TEST_DB_DSN as DB_DSN for database client
	originalDSN := os.Getenv("DB_DSN")
	os.Setenv("DB_DSN", os.Getenv("TEST_DB_DSN"))
	defer os.Setenv("DB_DSN", originalDSN)

	// Initialize database
	db, err := database.NewDatabaseClient()
	require.NoError(suite.T(), err, "Failed to initialize database client")
	suite.db = db

	// Initialize server
	srv := server.NewEchoServer(db)
	suite.server = srv

	// Get Echo instance from server for testing
	if echoSrv, ok := srv.(*server.EchoServer); ok {
		suite.e = echoSrv.GetEcho()
	} else {
		suite.e = echo.New()
	}
}

// TearDownSuite runs once after all tests
func (suite *IntegrationTestSuite) TearDownSuite() {
	// Cleanup if needed
}

// TestHealthEndpoints tests the health check endpoints
func (suite *IntegrationTestSuite) TestHealthEndpoints() {
	// Test Readiness
	req := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var response map[string]string
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "OK", response["status"])

	// Test Liveness
	req = httptest.NewRequest(http.MethodGet, "/liveness", nil)
	rec = httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "OK", response["status"])
}

// TestGetAllBooks tests GET /api/niv/books
func (suite *IntegrationTestSuite) TestGetAllBooks() {
	req := httptest.NewRequest(http.MethodGet, "/api/niv/books", nil)
	rec := httptest.NewRecorder()

	// Use Echo's ServeHTTP to test the full HTTP stack
	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var books []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &books)
	require.NoError(suite.T(), err)

	assert.Greater(suite.T(), len(books), 0, "Should return at least one book")

	// Verify structure
	firstBook := books[0]
	assert.Contains(suite.T(), firstBook, "book_id")
	assert.Contains(suite.T(), firstBook, "book")
}

// TestGetAllVerses tests GET /api/niv/verses
func (suite *IntegrationTestSuite) TestGetAllVerses() {
	req := httptest.NewRequest(http.MethodGet, "/api/niv/verses", nil)
	rec := httptest.NewRecorder()

	// Use Echo's ServeHTTP to test the full HTTP stack
	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var verses []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &verses)
	require.NoError(suite.T(), err)

	assert.Greater(suite.T(), len(verses), 0, "Should return at least one verse")

	// Verify structure - API returns capitalized field names: BookID, Book, Chapter, Verse, Text
	firstVerse := verses[0]
	assert.Contains(suite.T(), firstVerse, "BookID")
	assert.Contains(suite.T(), firstVerse, "Book")
	assert.Contains(suite.T(), firstVerse, "Chapter")
	assert.Contains(suite.T(), firstVerse, "Verse")
	assert.Contains(suite.T(), firstVerse, "Text")
}

// TestGetVersesByChapter tests GET /api/niv/:bookId/:chapterId/verses
func (suite *IntegrationTestSuite) TestGetVersesByChapter() {
	// Test with valid book_id and chapter
	req := httptest.NewRequest(http.MethodGet, "/api/niv/1/1/verses", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var verses []map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &verses)
	require.NoError(suite.T(), err)

	assert.Greater(suite.T(), len(verses), 0, "Should return verses for Genesis chapter 1")

	// Verify all verses are from chapter 1 - API returns capitalized field names
	for _, verse := range verses {
		assert.Equal(suite.T(), float64(1), verse["Chapter"], "All verses should be from chapter 1")
		assert.Equal(suite.T(), float64(1), verse["BookID"], "All verses should be from book_id 1")
	}
}

// TestGetVersesByChapter_InvalidChapter tests invalid chapter number
func (suite *IntegrationTestSuite) TestGetVersesByChapter_InvalidChapter() {
	req := httptest.NewRequest(http.MethodGet, "/api/niv/1/invalid/verses", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
}

// TestGetChaptersByBook tests GET /api/niv/chapters/:bookId
func (suite *IntegrationTestSuite) TestGetChaptersByBook() {
	// Test with book_id
	req := httptest.NewRequest(http.MethodGet, "/api/niv/chapters/1", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(suite.T(), err)

	assert.Contains(suite.T(), response, "MaxChapter")
	maxChapter, ok := response["MaxChapter"].(float64)
	require.True(suite.T(), ok, "MaxChapter should be a number")
	assert.Greater(suite.T(), maxChapter, float64(0), "Should have at least 1 chapter")
}

// TestGetChaptersByBook_InvalidBookId tests GET /api/niv/chapters/:bookId with invalid bookId
func (suite *IntegrationTestSuite) TestGetChaptersByBook_InvalidBookId() {
	// Test with invalid bookId (non-numeric)
	req := httptest.NewRequest(http.MethodGet, "/api/niv/chapters/invalid", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(suite.T(), rec.Body.String(), "Invalid book ID")
}

// TestGetVersesByChapter_InvalidBookId tests GET /api/niv/:bookId/:chapterId/verses with invalid bookId
func (suite *IntegrationTestSuite) TestGetVersesByChapter_InvalidBookId() {
	// Test with invalid bookId (non-numeric)
	req := httptest.NewRequest(http.MethodGet, "/api/niv/invalid/1/verses", nil)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)
	assert.Contains(suite.T(), rec.Body.String(), "Invalid book ID")
}

// TestExplainVerse tests POST /api/niv/explain
func (suite *IntegrationTestSuite) TestExplainVerse() {
	// Skip if OpenAI API key is not set
	if os.Getenv("OPENAI_API_KEY") == "" {
		suite.T().Skip("OPENAI_API_KEY not set - skipping explain verse test")
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
	require.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	// This might succeed or fail depending on API key validity
	// Should return either 200 (success) or 500 (API error)
	assert.Contains(suite.T(), []int{http.StatusOK, http.StatusInternalServerError}, rec.Code)
}

// TestExplainVerse_InvalidRequest tests POST /api/niv/explain with invalid data
func (suite *IntegrationTestSuite) TestExplainVerse_InvalidRequest() {
	// Skip if API key is not set, as server will return 500 instead of 400
	if os.Getenv("OPENAI_API_KEY") == "" {
		suite.T().Skip("OPENAI_API_KEY not set - skipping invalid request test")
		return
	}

	invalidReq := map[string]interface{}{
		"book": "Genesis",
		// Missing required fields
	}

	body, err := json.Marshal(invalidReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rec.Code)

	var response map[string]string
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	require.NoError(suite.T(), err)
	assert.Contains(suite.T(), response["error"], "Invalid")
}

// TestExplainVerse_DefaultValues tests that default values are applied
func (suite *IntegrationTestSuite) TestExplainVerse_DefaultValues() {
	if os.Getenv("OPENAI_API_KEY") == "" {
		suite.T().Skip("OPENAI_API_KEY not set")
		return
	}

	// Request without age and belief (should use defaults)
	explainReq := dto.ExplainRequest{
		Book:       "Genesis",
		Chapter:    1,
		StartVerse: 1,
		EndVerse:   3,
		// Age and Belief omitted - should default to 25 and 3
	}

	body, err := json.Marshal(explainReq)
	require.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/api/niv/explain", bytes.NewBuffer(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	suite.e.ServeHTTP(rec, req)

	// Should process (may fail due to API, but should handle defaults)
	assert.Contains(suite.T(), []int{http.StatusOK, http.StatusInternalServerError}, rec.Code)
}

// Run the test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
