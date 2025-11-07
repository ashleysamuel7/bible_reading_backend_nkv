#!/bin/bash

# Test script to verify the chapter/book changes
# This script tests the updated API endpoints

BASE_URL="${BASE_URL:-http://localhost:8000}"

echo "=========================================="
echo "Testing Chapter/Book API Changes"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

test_count=0
pass_count=0
fail_count=0

# Test function
test_endpoint() {
    local name="$1"
    local method="$2"
    local url="$3"
    local expected_status="$4"
    
    test_count=$((test_count + 1))
    echo -n "Test $test_count: $name ... "
    
    response=$(curl -s -w "\n%{http_code}" -X "$method" "$BASE_URL$url")
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -eq "$expected_status" ]; then
        echo -e "${GREEN}PASS${NC} (Status: $http_code)"
        pass_count=$((pass_count + 1))
        return 0
    else
        echo -e "${RED}FAIL${NC} (Expected: $expected_status, Got: $http_code)"
        echo "  Response: $body"
        fail_count=$((fail_count + 1))
        return 1
    fi
}

echo "Testing valid endpoints..."
echo ""

# Test 1: Get chapters for book 1 (Genesis)
test_endpoint "Get chapters for book 1" "GET" "/api/niv/chapters/1" 200

# Test 2: Get verses for book 1, chapter 1
test_endpoint "Get verses for book 1, chapter 1" "GET" "/api/niv/1/1/verses" 200

# Test 3: Get all books
test_endpoint "Get all books" "GET" "/api/niv/books" 200

echo ""
echo "Testing invalid input validation..."
echo ""

# Test 4: Invalid bookId (non-numeric) for chapters
test_endpoint "Invalid bookId for chapters (non-numeric)" "GET" "/api/niv/chapters/invalid" 400

# Test 5: Invalid bookId (non-numeric) for verses
test_endpoint "Invalid bookId for verses (non-numeric)" "GET" "/api/niv/invalid/1/verses" 400

# Test 6: Invalid chapterId (non-numeric) for verses
test_endpoint "Invalid chapterId for verses (non-numeric)" "GET" "/api/niv/1/invalid/verses" 400

# Test 7: Invalid bookId (negative number)
test_endpoint "Invalid bookId (negative)" "GET" "/api/niv/chapters/-1" 400

echo ""
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo "Total tests: $test_count"
echo -e "${GREEN}Passed: $pass_count${NC}"
if [ $fail_count -gt 0 ]; then
    echo -e "${RED}Failed: $fail_count${NC}"
else
    echo -e "${GREEN}Failed: $fail_count${NC}"
fi
echo ""

if [ $fail_count -eq 0 ]; then
    echo -e "${GREEN}All tests passed! ✓${NC}"
    exit 0
else
    echo -e "${RED}Some tests failed! ✗${NC}"
    exit 1
fi

