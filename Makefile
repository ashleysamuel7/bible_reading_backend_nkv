.PHONY: test test-integration test-coverage test-report clean help

# Test configuration
TEST_DIR = test
REPORTS_DIR = $(TEST_DIR)/reports
FIXTURES_DIR = $(TEST_DIR)/fixtures
COVERAGE_FILE = $(REPORTS_DIR)/coverage.out
COVERAGE_HTML = $(REPORTS_DIR)/coverage.html

# Ensure directories exist
$(shell mkdir -p $(REPORTS_DIR) $(FIXTURES_DIR))

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

test: ## Run all Go tests
	@echo "Running all Go tests..."
	@go test -v ./...

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	@go test -v ./test/integration/...

test-coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

test-report: test-coverage ## Generate test coverage report
	@echo "Test coverage report available at: $(COVERAGE_HTML)"

test-python: ## Run Python endpoint tests
	@echo "Running Python endpoint tests..."
	@python3 test_endpoints.py

test-all: test test-python ## Run all tests (Go + Python)

clean: ## Clean test artifacts
	@echo "Cleaning test artifacts..."
	@rm -f $(REPORTS_DIR)/*.out $(REPORTS_DIR)/*.html
	@rm -f $(REPORTS_DIR)/API_TEST_REPORT_*.md
	@echo "Clean complete"

clean-all: clean ## Clean all test artifacts including reports
	@rm -rf $(REPORTS_DIR)/*
	@echo "All test artifacts cleaned"

