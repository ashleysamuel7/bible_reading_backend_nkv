# Test Cases Execution Guide

The test cases in `test_cases.json` are test specifications. Here are several ways to execute them:

## Quick Start - Bash Script (Easiest)

**Prerequisites:**
- Backend server running on `http://localhost:8000`
- `curl` installed
- (Optional) `OPENAI_API_KEY` set for explain endpoint tests

**Run all tests:**
```bash
cd prompts/tests
./run_tests.sh
```

**Run against different server:**
```bash
./run_tests.sh http://13.203.234.131:8000
```

The script will:
- ✅ Run all basic endpoint tests
- ✅ Show pass/fail for each test
- ✅ Display summary at the end

## Option 1: Manual Testing with curl

Test individual endpoints:

```bash
# Health check
curl -X GET http://localhost:8000/readiness

# Get all verses
curl -X GET http://localhost:8000/api/niv/verses

# Get verses by book and chapter
curl -X GET http://localhost:8000/api/niv/1/1/verses

# Get all books
curl -X GET http://localhost:8000/api/niv/books

# Get max chapter
curl -X GET http://localhost:8000/api/niv/chapters/1

# Get explanation (requires OPENAI_API_KEY)
curl -X POST http://localhost:8000/api/niv/explain \
  -H "Content-Type: application/json" \
  -d '{
    "book": "Gen",
    "chapter": 1,
    "start_verse": 1,
    "end_verse": 3
  }'
```

## Option 2: Go Test Runner

**Prerequisites:**
```bash
go get github.com/stretchr/testify/assert
go get github.com/stretchr/testify/require
```

**Run tests:**
```bash
cd server
go test -v -run TestFromJSON
```

**Note:** The Go test file (`server_test.go`) is a template. You'll need to implement the actual HTTP request execution using Echo's test utilities.

## Option 3: Using Postman/Newman

1. Import the test cases into Postman manually, or
2. Use Newman CLI to run JSON tests:
```bash
npm install -g newman
# Convert test_cases.json to Postman collection format first
newman run test_cases.json
```

## Option 4: Python Test Runner

Create a Python script using `requests` library:

```python
import json
import requests

with open('test_cases.json') as f:
    tests = json.load(f)

base_url = "http://localhost:8000"

for endpoint_test in tests['tests']:
    for test_case in endpoint_test['test_cases']:
        req = test_case['request']
        url = req['url'].replace('http://localhost:8000', base_url)
        
        if req['method'] == 'GET':
            response = requests.get(url)
        elif req['method'] == 'POST':
            response = requests.post(url, json=req['body'])
        
        expected = test_case['expected_status_codes']
        if response.status_code in expected:
            print(f"✓ {test_case['id']}: PASSED")
        else:
            print(f"✗ {test_case['id']}: FAILED (got {response.status_code})")
```

## Test Execution Tips

1. **Start the backend server first:**
   ```bash
   cd bible_reading_backend_nkv
   go run main.go
   ```

2. **Ensure database is running:**
   ```bash
   # Check database connection
   mysql -h localhost -u user -p bible_reading_db
   ```

3. **Set environment variables:**
   ```bash
   export DB_DSN="user:password@tcp(localhost:3306)/bible_reading_db"
   export OPENAI_API_KEY="sk-proj-your-key-here"
   ```

4. **Run tests in Docker:**
   ```bash
   # Start backend in Docker
   docker run -d --name bible-backend \
     --network host \
     --env-file .env \
     bible-backend:latest
   
   # Run tests
   ./run_tests.sh http://localhost:8000
   ```

## Expected Results

- **Health endpoints**: Should always return 200 OK
- **NIV endpoints**: Return 200 OK if database has data, 500 if DB error
- **Explain endpoint**: Returns 200 OK if OpenAI API key is valid, 500 if missing/invalid

