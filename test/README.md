# Test Directory Structure

This directory contains all test files, fixtures, and reports for the Bible Reading Backend.

## Directory Structure

```
test/
├── integration/          # Go integration tests
│   ├── api_test.go      # API endpoint tests
│   ├── endpoints_test.go # Endpoint-specific tests
│   ├── helpers_test.go  # Test helper functions
│   ├── setup_test.go    # Test setup and teardown
│   └── README.md        # Integration test documentation
├── fixtures/            # Test fixtures and mock data
│   └── .gitkeep        # Keep directory in git
└── reports/             # Test reports and coverage
    ├── .gitkeep        # Keep directory in git
    ├── coverage.out    # Go test coverage (binary)
    ├── coverage.html   # Go test coverage (HTML)
    └── API_TEST_REPORT_*.md  # Python test reports
```

## Running Tests

### Go Tests

```bash
# Run all tests
make test

# Run integration tests only
make test-integration

# Run tests with coverage
make test-coverage

# Generate coverage report
make test-report
```

Or directly with Go:

```bash
# Run all tests
go test -v ./...

# Run integration tests
go test -v ./test/integration/...

# Run with coverage
go test -v -coverprofile=test/reports/coverage.out ./...
go tool cover -html=test/reports/coverage.out -o test/reports/coverage.html
```

### Python Tests

```bash
# Run Python endpoint tests
make test-python

# Or directly
python3 test_endpoints.py
```

### All Tests

```bash
# Run both Go and Python tests
make test-all
```

## Test Reports

- **Go Coverage Reports**: `test/reports/coverage.html`
- **Python Test Reports**: `test/reports/API_TEST_REPORT_YYYYMMDD_HHMMSS.md`

Reports are automatically generated when running tests and saved to the `test/reports/` directory.

## Test Fixtures

Place any test fixtures, mock data, or test resources in the `test/fixtures/` directory.

## Notes

- Test reports are gitignored but the directory structure is preserved
- Coverage reports are generated in HTML format for easy viewing
- Python test reports include timestamps in the filename
- All test artifacts are cleaned with `make clean`

