#!/bin/bash

# Simple test runner script for API endpoints
# Usage: ./run_tests.sh [base_url]
# Default base_url: http://localhost:8000

BASE_URL="${1:-http://localhost:8000}"
PASSED=0
FAILED=0

echo "🧪 Running API Tests against: $BASE_URL"
echo "=========================================="

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
run_test() {
    local test_id=$1
    local method=$2
    local endpoint=$3
    local expected_status=$4
    local data=$5
    local description=$6
    
    echo -n "Testing $test_id: $description ... "
    
    if [ -z "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$endpoint" \
            -H "Content-Type: application/json" \
            -d "$data")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}✓ PASSED${NC} (Status: $http_code)"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}✗ FAILED${NC} (Expected: $expected_status, Got: $http_code)"
        echo -e "${YELLOW}Response: $body${NC}"
        ((FAILED++))
        return 1
    fi
}

# Health Check Tests
echo ""
echo "📋 Health Check Tests"
echo "-------------------"
run_test "HEALTH.READINESS.001" "GET" "/readiness" 200 "" "Readiness check"
run_test "HEALTH.LIVENESS.001" "GET" "/liveness" 200 "" "Liveness check"

# NIV Verses Tests
echo ""
echo "📋 NIV Verses Tests"
echo "-----------------"
run_test "NIV.VERSES.001" "GET" "/api/niv/verses" 200 "" "Get all verses"
run_test "NIV.BOOKS.001" "GET" "/api/niv/books" 200 "" "Get all books"
run_test "NIV.CHAPTERS.001" "GET" "/api/niv/chapters/1" 200 "" "Get max chapter for book 1"
run_test "NIV.VERSES_BY_CHAPTER.001" "GET" "/api/niv/1/1/verses" 200 "" "Get verses for book 1, chapter 1"

# Explain Endpoint Tests
echo ""
echo "📋 Explain Endpoint Tests"
echo "-----------------------"

# Check if OPENAI_API_KEY is set
if [ -z "$OPENAI_API_KEY" ]; then
    echo -e "${YELLOW}⚠ Warning: OPENAI_API_KEY not set. Skipping explain tests.${NC}"
else
    explain_data='{"book":"Gen","chapter":1,"start_verse":1,"end_verse":3,"age":25,"belief":3}'
    run_test "NIV.EXPLAIN.001" "POST" "/api/niv/explain" 200 "$explain_data" "Get explanation with valid request"
    
    explain_data_minimal='{"book":"Gen","chapter":1,"start_verse":1,"end_verse":3}'
    run_test "NIV.EXPLAIN.005" "POST" "/api/niv/explain" 200 "$explain_data_minimal" "Get explanation with default age/belief"
    
    # Test invalid request
    invalid_data='{"book":"Gen"}'
    run_test "NIV.EXPLAIN.004" "POST" "/api/niv/explain" 400 "$invalid_data" "Get explanation with missing fields"
fi

# Summary
echo ""
echo "=========================================="
echo "📊 Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total: $((PASSED + FAILED))"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✅ All tests passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Some tests failed!${NC}"
    exit 1
fi

