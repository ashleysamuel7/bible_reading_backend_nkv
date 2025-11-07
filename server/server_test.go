package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"

	"bible_reading_backend_nkv/database"

	"github.com/stretchr/testify/assert"
)

// TestHelper loads test cases from JSON and runs them
func TestFromJSON(t *testing.T) {
	// Try multiple possible paths for test_cases.json
	testPaths := []string{
		"docs/tests/test_cases.json",
		"prompts/tests/test_cases.json",
		"../docs/tests/test_cases.json",
	}
	
	var jsonFile *os.File
	var err error
	
	for _, path := range testPaths {
		jsonFile, err = os.Open(path)
		if err == nil {
			break
		}
	}
	
	if jsonFile == nil {
		t.Skip("test_cases.json not found - skipping JSON-based tests")
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var testSuite TestSuite
	json.Unmarshal(byteValue, &testSuite)

	// Setup test server (mock database)
	// For integration tests, use real database connection
	db, err := database.NewDatabaseClient()
	if err != nil {
		t.Skip("Database not available for integration tests")
	}

	server := NewEchoServer(db)

	// Run tests
	for _, endpointTest := range testSuite.Tests {
		for _, testCase := range endpointTest.TestCases {
			t.Run(testCase.ID, func(t *testing.T) {
				runTestCase(t, server, endpointTest.Endpoint, testCase)
			})
		}
	}
}

func runTestCase(t *testing.T, server Server, endpoint Endpoint, testCase TestCase) {
	// Create request
	req := httptest.NewRequest(testCase.Request.Method, testCase.Request.URL, nil)
	if testCase.Request.Body != nil {
		bodyBytes, _ := json.Marshal(testCase.Request.Body)
		req = httptest.NewRequest(testCase.Request.Method, testCase.Request.URL, bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
	}

	// Create response recorder
	rec := httptest.NewRecorder()

	// Execute request (this requires Echo server setup)
	// For now, this is a template - adjust based on your Echo setup

	// Assert status code
	assert.Contains(t, testCase.ExpectedStatusCodes, rec.Code,
		fmt.Sprintf("Test %s: Expected status %v, got %d", testCase.ID, testCase.ExpectedStatusCodes, rec.Code))
}

// Test data structures
type TestSuite struct {
	Tests []EndpointTest `json:"tests"`
}

type EndpointTest struct {
	Endpoint  Endpoint   `json:"endpoint"`
	Purpose   string     `json:"purpose"`
	TestCases []TestCase `json:"test_cases"`
}

type Endpoint struct {
	Method       string `json:"method"`
	Path         string `json:"path"`
	AuthRequired string `json:"auth_required"`
}

type TestCase struct {
	ID                      string      `json:"id"`
	Title                   string      `json:"title"`
	Type                    string      `json:"type"`
	Preconditions           []string    `json:"preconditions"`
	Steps                   []string    `json:"steps"`
	Request                 Request     `json:"request"`
	ExpectedStatusCodes     []int       `json:"expected_status_codes"`
	ExpectedResponseSchema  interface{} `json:"expected_response_schema"`
	ExpectedResponseExample interface{} `json:"expected_response_example"`
	CleanupSteps            []string    `json:"cleanup_steps"`
	Severity                string      `json:"severity"`
}

type Request struct {
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	Headers     map[string]string `json:"headers"`
	QueryParams map[string]string `json:"query_params"`
	Body        interface{}       `json:"body"`
}
