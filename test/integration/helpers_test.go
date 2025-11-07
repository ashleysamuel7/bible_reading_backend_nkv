package integration

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"

	"bible_reading_backend_nkv/server"
)

// TestHelper provides helper functions for integration tests
type TestHelper struct {
	Server server.Server
}

// NewTestHelper creates a new test helper
func NewTestHelper(server server.Server) *TestHelper {
	return &TestHelper{
		Server: server,
	}
}

// MakeRequest makes an HTTP request and returns the response
func (h *TestHelper) MakeRequest(method, path string, body interface{}, headers map[string]string) (*httptest.ResponseRecorder, error) {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	if body != nil && headers["Content-Type"] == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	
	// This helper is not fully implemented yet
	// Tests should use echo.ServeHTTP directly
	return rec, nil
}

// ParseJSONResponse parses a JSON response into a map or slice
func ParseJSONResponse(rec *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(rec.Body.Bytes(), target)
}

