# API Reference Documentation

## Base URL
```
http://localhost:8000
```

## Authentication

All protected endpoints require JWT authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

## Endpoints

### Authentication

#### Register User
```http
POST /api/register/
Content-Type: application/json

{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "password": "password123",
  "age": 30,
  "belif_rating": 4
}
```

**Response:**
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
```http
POST /api/login/
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "access": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### User Profile

#### Get Current User
```http
GET /api/users/me
Authorization: Bearer <token>
```

**Response:**
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
```http
PUT /api/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "age": 31
}
```

**Response:** Updated user object (same format as GET /api/users/me)

#### Delete Current User
```http
DELETE /api/users/me
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

### Favorite Verses

#### Add Favorite Verse
```http
POST /api/users/me/favorites
Authorization: Bearer <token>
Content-Type: application/json

{
  "book_id": 1,
  "chapter": 1,
  "verse": 1
}
```

**Response:**
```json
{
  "message": "Favorite verse added successfully"
}
```

#### Get Favorite Verses
```http
GET /api/users/me/favorites?page=1&limit=20
Authorization: Bearer <token>
```

**Query Parameters:**
- `page` (optional): Page number, default: 1
- `limit` (optional): Items per page, default: 20, max: 100

**Response:**
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
```http
DELETE /api/users/me/favorites/1/1/1
Authorization: Bearer <token>
```

**Path Parameters:**
- `book_id`: Book ID
- `chapter`: Chapter number
- `verse`: Verse number

**Response:**
```json
{
  "message": "Favorite verse removed successfully"
}
```

### Highlighted Verses

#### Add Highlighted Verse
```http
POST /api/users/me/highlights
Authorization: Bearer <token>
Content-Type: application/json

{
  "book_id": 1,
  "chapter": 1,
  "verse": 1,
  "note": "Important verse about creation",
  "color": "yellow"
}
```

**Response:**
```json
{
  "message": "Highlighted verse added successfully"
}
```

#### Get Highlighted Verses
```http
GET /api/users/me/highlights?page=1&limit=20
Authorization: Bearer <token>
```

**Response:**
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
```http
PUT /api/users/me/highlights/1/1/1
Authorization: Bearer <token>
Content-Type: application/json

{
  "note": "Updated note",
  "color": "blue"
}
```

**Response:**
```json
{
  "message": "Highlighted verse updated successfully"
}
```

#### Remove Highlighted Verse
```http
DELETE /api/users/me/highlights/1/1/1
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Highlighted verse removed successfully"
}
```

### Last Read

#### Update Last Read Position
```http
POST /api/users/me/last-read
Authorization: Bearer <token>
Content-Type: application/json

{
  "book_id": 1,
  "book_name": "Genesis",
  "chapter": 1,
  "verse": 10
}
```

**Response:**
```json
{
  "message": "Last read updated successfully"
}
```

#### Get Last Read Position
```http
GET /api/users/me/last-read
Authorization: Bearer <token>
```

**Response:**
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

#### Get Last Read Verses (Legacy)
```http
GET /api/last-read-verses/
Authorization: Bearer <token>
```

**Response:**
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

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request parameters"
}
```

### 401 Unauthorized
```json
{
  "error": "Invalid or expired token"
}
```

### 404 Not Found
```json
{
  "error": "Verse not found"
}
```

### 409 Conflict
```json
{
  "error": "Verse already in favorites"
}
```

### 500 Internal Server Error
```json
{
  "error": "Failed to process request"
}
```

## Status Codes

- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `401 Unauthorized` - Authentication required or invalid token
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate entry (e.g., verse already favorited)
- `500 Internal Server Error` - Server error

