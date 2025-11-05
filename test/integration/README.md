# Integration Tests

This directory contains integration tests for the Bible Reading Backend API.

## Overview

Integration tests verify that all components work together correctly:
- API endpoints
- Database connections
- Business logic
- Error handling

## Test Structure

### Files

- `api_test.go` - Main test suite using testify/suite
- `endpoints_test.go` - Individual endpoint tests
- `helpers_test.go` - Helper functions for testing

## Running Tests

### Prerequisites

1. **Database Setup**: Set `TEST_DB_DSN` environment variable
   ```bash
   export TEST_DB_DSN="user:password@tcp(localhost:3306)/bible_db?charset=utf8mb4&parseTime=True&loc=UTC"
   ```

2. **OpenAI API Key** (optional, for explain endpoint tests):
   ```bash
   export OPENAI_API_KEY="sk-your-key-here"
   ```

### Run All Integration Tests

```bash
# From project root
go test ./test/integration -v
```

### Run Specific Test

```bash
go test ./test/integration -v -run TestEndpoints
```

### Run with Coverage

```bash
go test ./test/integration -v -cover
```

## Test Coverage

### Health Endpoints
- ✅ `GET /readiness` - Database connectivity check
- ✅ `GET /liveness` - Application health check

### NIV API Endpoints
- ✅ `GET /api/niv/books` - Get all books
- ✅ `GET /api/niv/verses` - Get all verses
- ✅ `GET /api/niv/:book/:chapter/verses` - Get verses by chapter
- ✅ `GET /api/niv/chapters/:book` - Get max chapters for book
- ✅ `POST /api/niv/explain` - Explain verses (requires OpenAI API key)

### Error Handling
- ✅ Invalid chapter numbers
- ✅ Missing required fields
- ✅ Missing API keys
- ✅ Invalid request formats

### Response Formats
- ✅ JSON structure validation
- ✅ Required fields presence
- ✅ Data type validation

## Test Environment

Tests use a real database connection (configured via `TEST_DB_DSN`). Make sure:

1. Test database is available
2. Database has the `niv` table populated
3. Database credentials are correct

## Skipping Tests

Tests automatically skip if:
- `TEST_DB_DSN` is not set
- `OPENAI_API_KEY` is not set (for explain endpoint tests)

## Continuous Integration

For CI/CD pipelines:

```yaml
# Example GitHub Actions
- name: Run Integration Tests
  env:
    TEST_DB_DSN: ${{ secrets.TEST_DB_DSN }}
    OPENAI_API_KEY: ${{ secrets.OPENAI_API_KEY }}
  run: go test ./test/integration -v
```

## Notes

- Tests use `httptest` for HTTP testing (no actual server required)
- Tests use real database connections (not mocked)
- OpenAI tests may fail if API key is invalid (expected behavior)
- Tests are designed to be idempotent (can run multiple times)

