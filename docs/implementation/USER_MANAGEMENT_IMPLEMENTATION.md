# User Management System Implementation Documentation

## Overview

This document describes the complete user management system implementation with JWT authentication, user profiles, and verse tracking features (favorites, highlights, and last read position).

## Table of Contents

1. [Database Schema](#database-schema)
2. [Authentication](#authentication)
3. [API Endpoints](#api-endpoints)
4. [Models and DTOs](#models-and-dtos)
5. [Database Layer](#database-layer)
6. [Server Implementation](#server-implementation)
7. [Setup and Configuration](#setup-and-configuration)
8. [Testing](#testing)
9. [Frontend Integration](#frontend-integration)

## Database Schema

### Users Table

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    age INT NOT NULL CHECK (age > 0 AND age < 150),
    believer_category TINYINT NOT NULL CHECK (believer_category >= 1 AND believer_category <= 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_email (email)
);
```

### User Favorite Verses Table

```sql
CREATE TABLE user_favorite_verses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    chapter INT NOT NULL,
    verse INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_favorite (user_id, book_id, chapter, verse),
    INDEX idx_user_id (user_id),
    INDEX idx_verse (book_id, chapter, verse)
);
```

### User Highlighted Verses Table

```sql
CREATE TABLE user_highlighted_verses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    book_id INT NOT NULL,
    chapter INT NOT NULL,
    verse INT NOT NULL,
    note TEXT,
    color VARCHAR(20) DEFAULT 'yellow',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_highlight (user_id, book_id, chapter, verse),
    INDEX idx_user_id (user_id),
    INDEX idx_verse (book_id, chapter, verse)
);
```

### User Last Read Table

```sql
CREATE TABLE user_last_read (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL UNIQUE,
    book_id INT NOT NULL,
    book_name VARCHAR(255) NOT NULL,
    chapter INT NOT NULL,
    verse INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id)
);
```

## Authentication

### JWT Token Authentication

The system uses JWT (JSON Web Tokens) for authentication. Tokens are generated upon login/registration and must be included in the `Authorization` header for protected endpoints.

**Token Format:**
```
Authorization: Bearer <jwt_token>
```

**Token Claims:**
- `user_id`: Integer user ID
- `exp`: Expiration timestamp
- `iat`: Issued at timestamp

**Configuration:**
- `JWT_SECRET`: Secret key for signing tokens (required)
- `JWT_EXPIRY_HOURS`: Token expiration time in hours (default: 24)

### Password Security

- Passwords are hashed using bcrypt with default cost
- Plain text passwords are never stored
- Minimum password length: 6 characters

## API Endpoints

### Authentication Endpoints (Public)

#### Register User
```
POST /api/register/
```

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30,
  "belif_rating": 4
}
```

**Response (201 Created):**
```json
{
  "access": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "john@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "name": "John Doe",
    "age": 30,
    "believer_category": 4,
    "created_at": "2024-11-05T14:00:00Z",
    "updated_at": "2024-11-05T14:00:00Z"
  }
}
```

#### Login
```
POST /api/login/
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "access": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### User Profile Endpoints (Protected - JWT Required)

#### Get Current User
```
GET /api/users/me
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "john@example.com",
  "first_name": "John",
  "last_name": "Doe",
  "name": "John Doe",
  "age": 30,
  "believer_category": 4,
  "created_at": "2024-11-05T14:00:00Z",
  "updated_at": "2024-11-05T14:00:00Z"
}
```

#### Update Current User
```
PUT /api/users/me
```

**Request Body (all fields optional):**
```json
{
  "first_name": "Jane",
  "last_name": "Smith",
  "age": 31,
  "believer_category": 5
}
```

**Response (200 OK):** Updated user object

#### Delete Current User
```
DELETE /api/users/me
```

**Response (200 OK):**
```json
{
  "message": "User deleted successfully"
}
```

### Favorite Verses Endpoints (Protected - JWT Required)

#### Add Favorite Verse
```
POST /api/users/me/favorites
```

**Request Body:**
```json
{
  "book_id": 1,
  "chapter": 1,
  "verse": 1
}
```

**Response (201 Created):**
```json
{
  "message": "Favorite verse added successfully"
}
```

#### Get Favorite Verses
```
GET /api/users/me/favorites?page=1&limit=20
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "book_id": 1,
      "book_name": "Genesis",
      "chapter": 1,
      "verse": 1,
      "text": "In the beginning God created the heavens and the earth.",
      "created_at": "2024-11-05T14:00:00Z"
    }
  ],
  "total": 10,
  "page": 1,
  "limit": 20,
  "total_pages": 1
}
```

#### Remove Favorite Verse
```
DELETE /api/users/me/favorites/:book_id/:chapter/:verse
```

**Response (200 OK):**
```json
{
  "message": "Favorite verse removed successfully"
}
```

### Highlighted Verses Endpoints (Protected - JWT Required)

#### Add Highlighted Verse
```
POST /api/users/me/highlights
```

**Request Body:**
```json
{
  "book_id": 1,
  "chapter": 1,
  "verse": 1,
  "note": "Important verse about creation",
  "color": "yellow"
}
```

**Response (201 Created):**
```json
{
  "message": "Highlighted verse added successfully"
}
```

#### Get Highlighted Verses
```
GET /api/users/me/highlights?page=1&limit=20
```

**Response (200 OK):**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "book_id": 1,
      "book_name": "Genesis",
      "chapter": 1,
      "verse": 1,
      "text": "In the beginning God created the heavens and the earth.",
      "note": "Important verse about creation",
      "color": "yellow",
      "created_at": "2024-11-05T14:00:00Z",
      "updated_at": "2024-11-05T14:00:00Z"
    }
  ],
  "total": 5,
  "page": 1,
  "limit": 20,
  "total_pages": 1
}
```

#### Update Highlighted Verse
```
PUT /api/users/me/highlights/:book_id/:chapter/:verse
```

**Request Body:**
```json
{
  "note": "Updated note",
  "color": "blue"
}
```

**Response (200 OK):**
```json
{
  "message": "Highlighted verse updated successfully"
}
```

#### Remove Highlighted Verse
```
DELETE /api/users/me/highlights/:book_id/:chapter/:verse
```

**Response (200 OK):**
```json
{
  "message": "Highlighted verse removed successfully"
}
```

### Last Read Endpoints (Protected - JWT Required)

#### Update Last Read Position
```
POST /api/users/me/last-read
```

**Request Body:**
```json
{
  "book_id": 1,
  "book_name": "Genesis",
  "chapter": 1,
  "verse": 10
}
```

**Response (200 OK):**
```json
{
  "message": "Last read updated successfully"
}
```

#### Get Last Read Position
```
GET /api/users/me/last-read
```

**Response (200 OK):**
```json
{
  "user_id": 1,
  "book_id": 1,
  "book_name": "Genesis",
  "chapter": 1,
  "verse": 10,
  "text": "And God called the dry land Earth...",
  "updated_at": "2024-11-05T14:00:00Z"
}
```

#### Get Last Read Verses (Legacy Endpoint)
```
GET /api/last-read-verses/
```

**Response (200 OK):**
```json
{
  "last_read_verses": [
    {
      "book_id": 1,
      "book_name": "Genesis",
      "chapter": 1,
      "verse": 10,
      "text": "And God called the dry land Earth..."
    }
  ]
}
```

## Models and DTOs

### Models

#### User Model (`models/user.go`)
```go
type User struct {
    ID              int       `gorm:"column:id;primaryKey;autoIncrement"`
    Email           string    `gorm:"column:email;uniqueIndex;not null"`
    Password        string    `gorm:"column:password;not null" json:"-"`
    FirstName       string    `gorm:"column:first_name;not null"`
    LastName        string    `gorm:"column:last_name;not null"`
    Age             int       `gorm:"column:age;not null"`
    BelieverCategory int      `gorm:"column:believer_category;not null"`
    CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
```

#### UserFavoriteVerse Model (`models/user_verse.go`)
```go
type UserFavoriteVerse struct {
    ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
    UserID    int       `gorm:"column:user_id;not null;index"`
    BookID    int       `gorm:"column:book_id;not null;index:idx_verse"`
    Chapter   int       `gorm:"column:chapter;not null;index:idx_verse"`
    Verse     int       `gorm:"column:verse;not null;index:idx_verse"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}
```

#### UserHighlightedVerse Model
```go
type UserHighlightedVerse struct {
    ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
    UserID    int       `gorm:"column:user_id;not null;index"`
    BookID    int       `gorm:"column:book_id;not null;index:idx_verse"`
    Chapter   int       `gorm:"column:chapter;not null;index:idx_verse"`
    Verse     int       `gorm:"column:verse;not null;index:idx_verse"`
    Note      string    `gorm:"column:note;type:text"`
    Color     string    `gorm:"column:color;size:20;default:yellow"`
    CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
```

#### UserLastRead Model
```go
type UserLastRead struct {
    ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
    UserID    int       `gorm:"column:user_id;uniqueIndex;not null"`
    BookID    int       `gorm:"column:book_id;not null"`
    BookName  string    `gorm:"column:book_name;not null;size:255"`
    Chapter   int       `gorm:"column:chapter;not null"`
    Verse     int       `gorm:"column:verse;not null"`
    UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}
```

## Database Layer

### User Management Methods

Located in `database/user.go`:

- `CreateUser(ctx, user)`: Creates a new user with hashed password
- `GetUserByID(ctx, id)`: Retrieves user by ID
- `GetUserByEmail(ctx, email)`: Retrieves user by email
- `UpdateUser(ctx, id, updates)`: Updates user fields
- `DeleteUser(ctx, id)`: Deletes user (cascades to related tables)
- `VerifyPassword(ctx, email, password)`: Verifies credentials and returns user

### Verse Tracking Methods

Located in `database/user_verse.go`:

**Favorite Verses:**
- `AddFavoriteVerse(ctx, userID, bookID, chapter, verse)`
- `GetFavoriteVerses(ctx, userID, limit, offset)`
- `GetFavoriteVersesCount(ctx, userID)`
- `RemoveFavoriteVerse(ctx, userID, bookID, chapter, verse)`
- `IsFavoriteVerse(ctx, userID, bookID, chapter, verse)`

**Highlighted Verses:**
- `AddHighlightedVerse(ctx, userID, bookID, chapter, verse, note, color)`
- `GetHighlightedVerses(ctx, userID, limit, offset)`
- `GetHighlightedVersesCount(ctx, userID)`
- `UpdateHighlightedVerse(ctx, userID, bookID, chapter, verse, note, color)`
- `RemoveHighlightedVerse(ctx, userID, bookID, chapter, verse)`

**Last Read:**
- `UpdateLastRead(ctx, userID, bookID, bookName, chapter, verse)`
- `GetLastRead(ctx, userID)`
- `GetLastReadVerses(ctx, userID)`

## Server Implementation

### Authentication Handlers (`server/auth_server.go`)

- `Register(ctx)`: Handles user registration
- `Login(ctx)`: Handles user login

### User Handlers (`server/user_server.go`)

- `GetCurrentUser(ctx)`: Returns current user profile
- `UpdateCurrentUser(ctx)`: Updates current user profile
- `DeleteCurrentUser(ctx)`: Deletes current user account

### Verse Tracking Handlers (`server/user_verse_server.go`)

- `AddFavoriteVerse(ctx)`: Adds verse to favorites
- `GetFavoriteVerses(ctx)`: Retrieves paginated favorite verses
- `RemoveFavoriteVerse(ctx)`: Removes verse from favorites
- `AddHighlightedVerse(ctx)`: Adds highlighted verse
- `GetHighlightedVerses(ctx)`: Retrieves paginated highlighted verses
- `UpdateHighlightedVerse(ctx)`: Updates highlight note/color
- `RemoveHighlightedVerse(ctx)`: Removes highlight
- `UpdateLastRead(ctx)`: Updates last read position
- `GetLastRead(ctx)`: Gets last read position
- `GetLastReadVerses(ctx)`: Legacy endpoint for frontend compatibility

### JWT Middleware (`server/middleware/auth.go`)

- Extracts and validates JWT tokens from Authorization header
- Stores user_id in Echo context for use in handlers
- Returns 401 Unauthorized for invalid/missing tokens

### JWT Utilities (`server/utils/jwt.go`)

- `GenerateToken(userID)`: Creates JWT token with user_id claim
- `ValidateToken(tokenString)`: Validates token and returns user_id

## Setup and Configuration

### Environment Variables

Create a `.env` file with the following variables:

```env
# Database Configuration
DB_DSN=user:password@tcp(localhost:3306)/database_name?charset=utf8mb4&parseTime=True&loc=Local

# JWT Configuration
JWT_SECRET=your-secret-key-here-change-in-production
JWT_EXPIRY_HOURS=24

# OpenAI API Key (for verse explanations)
OPENAI_API_KEY=your-openai-api-key-here
```

### Database Migration

The application automatically runs migrations on startup. Tables are created using GORM AutoMigrate:

```go
db.AutoMigrate(
    &models.User{},
    &models.UserFavoriteVerse{},
    &models.UserHighlightedVerse{},
    &models.UserLastRead{},
)
```

### Running the Server

```bash
cd bible_reading_backend_nkv
go run main.go
```

The server will start on port 8000 by default.

## Testing

### Python Integration Tests

Run the Python test suite:

```bash
python3 test_endpoints.py
```

Test reports are saved to `test/reports/API_TEST_REPORT_*.md`

### Go Integration Tests

Run Go tests:

```bash
# Set test database connection
export TEST_DB_DSN="user:password@tcp(localhost:3306)/test_db"

# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=test/reports/coverage.out ./...
go tool cover -html=test/reports/coverage.out -o test/reports/coverage.html
```

### Using Makefile

```bash
# Run all tests
make test-all

# Run Go tests only
make test

# Run Python tests only
make test-python

# Run with coverage
make test-coverage
```

## Frontend Integration

### Authentication Flow

1. User registers/logs in via `/api/register/` or `/api/login/`
2. Frontend receives JWT token in `access` field
3. Store token in cookie/localStorage
4. Include token in Authorization header for protected endpoints:
   ```
   Authorization: Bearer <token>
   ```

### Example Frontend API Calls

#### Register User
```javascript
const response = await fetch('http://localhost:8000/api/register/', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    first_name: 'John',
    last_name: 'Doe',
    email: 'john@example.com',
    password: 'password123',
    age: 30,
    belif_rating: 4
  })
});
const { access, user } = await response.json();
// Store token: localStorage.setItem('token', access);
```

#### Add Favorite Verse
```javascript
const token = localStorage.getItem('token');
const response = await fetch('http://localhost:8000/api/users/me/favorites', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    book_id: 1,
    chapter: 1,
    verse: 1
  })
});
```

#### Get Favorite Verses (with pagination)
```javascript
const token = localStorage.getItem('token');
const response = await fetch('http://localhost:8000/api/users/me/favorites?page=1&limit=20', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
const { data, total, page, limit, total_pages } = await response.json();
```

## Error Responses

All endpoints return consistent error formats:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `200 OK`: Success
- `201 Created`: Resource created
- `400 Bad Request`: Invalid request parameters
- `401 Unauthorized`: Invalid or missing token
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate entry (e.g., verse already favorited)
- `500 Internal Server Error`: Server error

## Validation Rules

- **Email**: Must be valid email format, unique
- **Password**: Minimum 6 characters, required
- **First Name**: Required, max 255 characters
- **Last Name**: Required, max 255 characters
- **Age**: Integer, 1-150
- **Believer Category**: Integer, 1-5
- **Verse References**: Must exist in NIV table (book_id, chapter, verse)
- **Favorite/Highlight**: One entry per user per verse (unique constraint)
- **Last Read**: One entry per user (unique constraint on user_id)

## Data Relationships

- Users can have many favorite verses (one-to-many)
- Users can have many highlighted verses (one-to-many)
- Users have one last read position (one-to-one)
- All verse references link to NIV table via composite key (book_id, chapter, verse)
- Foreign key constraints ensure data integrity with CASCADE delete

## Security Considerations

1. **Password Security**: All passwords are hashed with bcrypt before storage
2. **JWT Tokens**: Tokens are signed with a secret key and include expiration
3. **Token Validation**: All protected endpoints validate tokens before processing
4. **SQL Injection**: GORM parameterized queries prevent SQL injection
5. **CORS**: Configured to allow specific origins only
6. **Input Validation**: All user inputs are validated before processing

## Performance Considerations

1. **Database Indexes**: All foreign keys and commonly queried fields are indexed
2. **Pagination**: List endpoints support pagination to limit response size
3. **Efficient Queries**: Verse text is joined from NIV table only when needed
4. **Connection Pooling**: GORM handles database connection pooling

## Future Enhancements

Potential improvements for future versions:

1. Refresh token mechanism for extended sessions
2. Password reset functionality
3. Email verification
4. Rate limiting for API endpoints
5. Bulk operations for favorites/highlights
6. Verse collections/reading plans
7. Sharing functionality for highlighted verses
8. Export functionality for user data

