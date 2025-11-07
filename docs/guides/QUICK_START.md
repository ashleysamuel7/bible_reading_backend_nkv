# Quick Start Guide - User Management System

## Prerequisites

- Go 1.23+ installed
- MySQL database running
- Environment variables configured

## Setup Steps

### 1. Install Dependencies

```bash
cd bible_reading_backend_nkv
go mod download
```

### 2. Configure Environment

Create a `.env` file in the project root:

```env
DB_DSN=user:password@tcp(localhost:3306)/bible_db?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY_HOURS=24
OPENAI_API_KEY=your-openai-api-key
```

### 3. Run Database Migrations

Migrations run automatically on server startup. Ensure your database is accessible.

### 4. Start the Server

```bash
go run main.go
```

The server will start on `http://localhost:8000`

## Testing the API

### 1. Register a User

```bash
curl -X POST http://localhost:8000/api/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "password123",
    "age": 30,
    "belif_rating": 4
  }'
```

**Response:**
```json
{
  "access": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {...}
}
```

Save the `access` token for subsequent requests.

### 2. Login (Alternative)

```bash
curl -X POST http://localhost:8000/api/login/ \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### 3. Get User Profile

```bash
curl -X GET http://localhost:8000/api/users/me \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 4. Add Favorite Verse

```bash
curl -X POST http://localhost:8000/api/users/me/favorites \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1,
    "chapter": 1,
    "verse": 1
  }'
```

### 5. Get Favorite Verses

```bash
curl -X GET "http://localhost:8000/api/users/me/favorites?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### 6. Add Highlighted Verse

```bash
curl -X POST http://localhost:8000/api/users/me/highlights \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1,
    "chapter": 1,
    "verse": 1,
    "note": "Important verse",
    "color": "yellow"
  }'
```

### 7. Update Last Read Position

```bash
curl -X POST http://localhost:8000/api/users/me/last-read \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1,
    "book_name": "Genesis",
    "chapter": 1,
    "verse": 10
  }'
```

## Running Tests

### Python Tests

```bash
python3 test_endpoints.py
```

Reports are saved to `test/reports/API_TEST_REPORT_*.md`

### Go Tests

```bash
# Set test database
export TEST_DB_DSN="user:password@tcp(localhost:3306)/test_db"

# Run tests
go test ./...

# Run with coverage
make test-coverage
```

## Using Postman

Import the Postman collection:
- File: `bible reading Go.postman_collection.json`
- Set environment variable `token` after login
- All protected endpoints will use the token automatically

## Frontend Integration

The frontend can use these endpoints directly. The API is designed to work with the existing React frontend:

1. Store JWT token in localStorage/cookie after login
2. Include token in Authorization header for protected endpoints
3. Use pagination for list endpoints
4. Handle error responses consistently

Example frontend code:
```javascript
const token = localStorage.getItem('token');
const response = await fetch('http://localhost:8000/api/users/me/favorites', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

## Troubleshooting

### Server won't start
- Check database connection string in `.env`
- Ensure MySQL is running
- Verify JWT_SECRET is set

### Authentication fails
- Check token is valid and not expired
- Verify token is included in Authorization header
- Ensure JWT_SECRET matches between server restarts

### Database errors
- Verify database exists
- Check user permissions
- Ensure AutoMigrate completes successfully

### Tests fail
- Set TEST_DB_DSN environment variable
- Ensure test database is accessible
- Check database connection settings

## Next Steps

1. Review `USER_MANAGEMENT_IMPLEMENTATION.md` for detailed documentation
2. Check `API_REFERENCE.md` for complete API documentation
3. Run integration tests to verify all endpoints
4. Integrate with frontend application
5. Configure production environment variables

