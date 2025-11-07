# Bible Reading Backend - Project Structure & Migration Guide

## 📋 Table of Contents
1. [Current Structure Analysis](#current-structure-analysis)
2. [Proposed Go Best Practices Structure](#proposed-go-best-practices-structure)
3. [Directory & File Explanations](#directory--file-explanations)
4. [Migration Roadmap](#migration-roadmap)

---

## 🔍 Current Structure Analysis

### Current Layout
```
bible_reading_backend_nkv/
├── main.go                 # Entry point (root level)
├── go.mod
├── go.sum
├── Dockerfile
├── server/                 # HTTP handlers + server setup (mixed concerns)
│   ├── server.go
│   ├── niv_server.go
│   └── server_test.go
├── database/              # DB client + queries (mixed)
│   ├── client.go
│   ├── niv.go
│   └── commonDTO.go
├── models/                # GORM models
│   ├── models.go
│   └── niv.go
├── dto/                   # DTOs for requests/responses
│   └── ExplainRequest.go
└── prompts/               # Markdown documentation
    ├── details/
    └── tests/
```

### Issues Identified
1. ❌ **No `cmd/` directory**: Entry point at root mixes concerns
2. ❌ **Mixed responsibilities**: Handlers, business logic, and DB access in same files
3. ❌ **Configuration scattered**: Environment loading in database client
4. ❌ **No service layer**: Business logic mixed with handlers
5. ❌ **No repository layer**: DB queries directly in client
6. ❌ **Tests co-located**: Test files mixed with implementation
7. ❌ **No `internal/` or `pkg/`**: No clear public/private boundaries
8. ❌ **No `configs/`**: Configuration management not centralized
9. ❌ **No `scripts/`**: Build/deployment scripts scattered
10. ❌ **Hardcoded values**: CORS origins, port numbers in code

---

## 🏗️ Proposed Go Best Practices Structure

### Recommended Layout
```
bible_reading_backend_nkv/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── api/
│   │   ├── handlers/                  # HTTP handlers (presentation layer)
│   │   │   ├── health.go
│   │   │   ├── niv_handler.go
│   │   │   └── handler_test.go
│   │   ├── middleware/                # HTTP middleware
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   ├── recovery.go
│   │   │   └── middleware_test.go
│   │   └── routes/                    # Route definitions
│   │       └── routes.go
│   │
│   ├── service/                       # Business logic layer
│   │   ├── niv_service.go
│   │   ├── explain_service.go
│   │   └── service_test.go
│   │
│   ├── repository/                    # Data access layer
│   │   ├── interfaces.go              # Repository interfaces
│   │   ├── niv_repository.go
│   │   └── repository_test.go
│   │
│   ├── config/                        # Configuration management
│   │   ├── config.go                  # Config struct and loader
│   │   ├── database.go                # DB-specific config
│   │   └── server.go                  # Server-specific config
│   │
│   └── database/                      # Database connection & setup
│       ├── client.go                  # GORM client wrapper
│       └── migrations.go              # DB migrations (optional)
│
├── pkg/                               # Public packages (if needed)
│   └── errors/                        # Shared error types
│       └── errors.go
│
├── models/                            # Domain models (GORM models)
│   ├── health.go
│   └── niv.go
│
├── dto/                               # Data Transfer Objects
│   ├── request/
│   │   └── explain_request.go
│   └── response/
│       └── explain_response.go
│
├── configs/                           # Configuration files
│   ├── .env.example                   # Example environment variables
│   └── config.yaml                    # Optional YAML config (if needed)
│
├── scripts/                           # Build & deployment scripts
│   ├── build.sh
│   ├── test.sh
│   ├── migrate.sh                     # Database migration script
│   └── docker-build.sh
│
├── .github/
│   └── workflows/
│       ├── ci.yml                     # CI/CD pipeline
│       └── docker.yml                 # Docker build workflow
│
├── test/                              # Integration tests
│   ├── integration/
│   │   └── api_test.go
│   └── fixtures/
│       └── test_data.sql
│
├── docs/                              # Documentation
│   ├── API.md                         # API documentation
│   └── DEPLOYMENT.md                  # Deployment guide
│
├── Dockerfile
├── docker-compose.yml                 # Local development setup
├── .dockerignore
├── .gitignore
├── go.mod
├── go.sum
├── Makefile                           # Common build commands
└── README.md
```

---

## 📁 Directory & File Explanations

### `cmd/` - Application Entry Points

**Purpose**: Contains main applications/executables for your project.

**Structure**:
- `cmd/api/main.go`: The main entry point that initializes and runs the HTTP server

**Why**: 
- Separates entry points from library code
- Allows multiple binaries (e.g., `cmd/api`, `cmd/migrate`, `cmd/worker`)
- Follows Go standard project layout

**Example**:
```go
// cmd/api/main.go
package main

import (
    "bible_reading_backend_nkv/internal/config"
    "bible_reading_backend_nkv/internal/database"
    "bible_reading_backend_nkv/internal/api"
)

func main() {
    cfg := config.Load()
    db := database.NewClient(cfg.Database)
    server := api.NewServer(cfg.Server, db)
    server.Start()
}
```

---

### `internal/` - Private Application Code

**Purpose**: Contains code that is private to your application and not meant to be imported by other projects.

**Why**: 
- Prevents external packages from importing your internal code
- Enforces encapsulation
- Standard Go practice for application code

#### `internal/api/` - HTTP Layer

**`handlers/`**: HTTP request handlers
- **Purpose**: Handle HTTP requests, validate input, call services
- **Responsibilities**:
  - Parse HTTP requests
  - Validate input (DTOs)
  - Call service layer
  - Format HTTP responses
  - Handle HTTP-specific errors
- **Example**: `niv_handler.go` handles `/api/niv/*` routes

**`middleware/`**: HTTP middleware
- **Purpose**: Cross-cutting concerns (CORS, logging, auth, recovery)
- **Files**:
  - `cors.go`: CORS configuration
  - `logger.go`: Request logging
  - `recovery.go`: Panic recovery
  - `auth.go`: Authentication (if needed)

**`routes/`**: Route definitions
- **Purpose**: Centralized route registration
- **Why**: Separates routing logic from server setup

#### `internal/service/` - Business Logic Layer

**Purpose**: Contains business logic and orchestration.

**Responsibilities**:
- Implement business rules
- Coordinate between repositories
- Call external APIs (e.g., OpenAI)
- Transform data between layers
- Handle business-level errors

**Example**:
```go
// internal/service/explain_service.go
type ExplainService interface {
    ExplainVerse(ctx context.Context, req dto.ExplainRequest) (string, error)
}

type explainService struct {
    openaiClient *openai.Client
}

func (s *explainService) ExplainVerse(ctx context.Context, req dto.ExplainRequest) (string, error) {
    // Business logic here
    // Call OpenAI API
    // Return explanation
}
```

**Why**: 
- Separates business logic from HTTP concerns
- Makes business logic testable without HTTP layer
- Allows reuse of business logic in other contexts (CLI, gRPC, etc.)

#### `internal/repository/` - Data Access Layer

**Purpose**: Encapsulates all database operations.

**Responsibilities**:
- Execute database queries
- Map database results to domain models
- Handle database-specific errors
- Provide clean interface for data access

**Example**:
```go
// internal/repository/interfaces.go
type NIVRepository interface {
    GetAllVerses(ctx context.Context) ([]models.NIV, error)
    GetVersesByChapter(ctx context.Context, book string, chapter int) ([]models.NIV, error)
    GetAllBooks(ctx context.Context) ([]dto.BookDTO, error)
    GetChaptersByBook(ctx context.Context, book string) (dto.ChapterMaxDTO, error)
}

// internal/repository/niv_repository.go
type nivRepository struct {
    db *gorm.DB
}

func (r *nivRepository) GetAllVerses(ctx context.Context) ([]models.NIV, error) {
    var verses []models.NIV
    err := r.db.WithContext(ctx).Find(&verses).Error
    return verses, err
}
```

**Why**:
- Separates data access from business logic
- Makes database layer easily testable with mocks
- Allows switching databases without changing business logic
- Centralizes SQL/GORM queries

#### `internal/config/` - Configuration Management

**Purpose**: Centralized configuration loading and validation.

**Structure**:
- `config.go`: Main config struct and loader
- `database.go`: Database configuration
- `server.go`: Server configuration (port, CORS, etc.)

**Example**:
```go
// internal/config/config.go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    OpenAI   OpenAIConfig
}

type ServerConfig struct {
    Port         string
    CORSOrigins  []string
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}

func Load() *Config {
    // Load from environment variables
    // Validate required fields
    // Return config
}
```

**Why**:
- Single source of truth for configuration
- Type-safe configuration
- Environment-specific configs
- Easy to test with different configs

#### `internal/database/` - Database Connection

**Purpose**: Manages database connection and initialization.

**Responsibilities**:
- Create and configure GORM client
- Handle connection pooling
- Provide health check methods
- Database migration helpers (optional)

**Why**: Separates connection management from repository logic

---

### `pkg/` - Public Packages

**Purpose**: Contains code that can be imported by external projects.

**When to use**: 
- Shared utilities
- Public APIs
- Reusable components

**Example**: `pkg/errors/` for custom error types that might be used by other services

**Note**: For this project, `pkg/` might be minimal or empty since it's an application, not a library.

---

### `models/` - Domain Models

**Purpose**: GORM models representing database entities.

**Structure**:
- `health.go`: Health check model
- `niv.go`: NIV Bible verses model

**Why**: 
- Domain models are at the root level (or in `pkg/models/`)
- They represent core business entities
- Used across layers (repository, service, handlers)

---

### `dto/` - Data Transfer Objects

**Purpose**: Structures for API requests and responses.

**Structure**:
- `request/`: Incoming request DTOs
  - `explain_request.go`
- `response/`: Outgoing response DTOs
  - `explain_response.go`
  - `health_response.go`

**Why**:
- Separates API contracts from domain models
- Allows different request/response formats
- Makes API versioning easier

---

### `configs/` - Configuration Files

**Purpose**: Static configuration files (not code).

**Files**:
- `.env.example`: Template for environment variables
- `config.yaml`: Optional YAML configuration (if needed)

**Why**: 
- Documents required environment variables
- Provides templates for different environments
- Version controlled (excluding actual `.env`)

---

### `scripts/` - Build & Deployment Scripts

**Purpose**: Automation scripts for common tasks.

**Files**:
- `build.sh`: Build the application
- `test.sh`: Run tests
- `migrate.sh`: Run database migrations
- `docker-build.sh`: Build Docker image

**Why**: 
- Standardizes build/deploy processes
- Makes CI/CD easier
- Documents common operations

---

### `test/` - Integration Tests

**Purpose**: Integration and end-to-end tests.

**Structure**:
- `integration/`: Full integration tests
- `fixtures/`: Test data

**Why**: 
- Separates unit tests (in `*_test.go` files) from integration tests
- Allows running integration tests separately
- Can include database setup/teardown

---

### `.github/workflows/` - CI/CD

**Purpose**: GitHub Actions workflows for automation.

**Files**:
- `ci.yml`: Continuous Integration (tests, linting)
- `docker.yml`: Docker build and push

**Why**: 
- Automated testing on PRs
- Automated builds and deployments
- Consistent development workflow

---

### Root Files

**`Makefile`**: Common commands
```makefile
.PHONY: build test run migrate docker-build

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

run:
	go run ./cmd/api

migrate:
	./scripts/migrate.sh

docker-build:
	docker build -t bible-backend .
```

**`docker-compose.yml`**: Local development
- MySQL database
- Application
- Environment variables

**`.dockerignore`**: Exclude files from Docker builds
- `.git/`, `test/`, `docs/`, etc.

---

## 🗺️ Migration Roadmap

### Phase 1: Foundation (Week 1)
**Goal**: Set up directory structure and configuration

#### Tasks:
1. ✅ Create directory structure
   ```bash
   mkdir -p cmd/api
   mkdir -p internal/{api/{handlers,middleware,routes},service,repository,config,database}
   mkdir -p pkg/errors
   mkdir -p dto/{request,response}
   mkdir -p configs
   mkdir -p scripts
   mkdir -p test/integration
   mkdir -p .github/workflows
   ```

2. ✅ Move and refactor configuration
   - Create `internal/config/config.go`
   - Move environment variable loading from `database/client.go`
   - Extract hardcoded values (CORS origins, port, etc.)
   - Create `configs/.env.example`

3. ✅ Create `cmd/api/main.go`
   - Move initialization logic from root `main.go`
   - Use new config package
   - Initialize dependencies in correct order

4. ✅ Update `go.mod` imports
   - Update import paths to match new structure

**Deliverables**:
- ✅ New directory structure created
- ✅ Configuration centralized
- ✅ Entry point moved to `cmd/api/main.go`
- ✅ Project compiles (may have temporary placeholders)

---

### Phase 2: Database & Repository Layer (Week 1-2)
**Goal**: Separate data access logic

#### Tasks:
1. ✅ Refactor `internal/database/`
   - Move `database/client.go` to `internal/database/client.go`
   - Keep only connection management
   - Remove query logic

2. ✅ Create repository layer
   - Create `internal/repository/interfaces.go` with interfaces
   - Move query logic from `database/niv.go` to `internal/repository/niv_repository.go`
   - Move DTOs from `database/commonDTO.go` to appropriate location
   - Implement repository pattern

3. ✅ Update tests
   - Move `database` tests to `internal/repository/repository_test.go`
   - Create mock repositories for testing

**Deliverables**:
- ✅ Repository layer implemented
- ✅ Database client only handles connections
- ✅ All DB queries moved to repositories
- ✅ Repository tests passing

---

### Phase 3: Service Layer (Week 2)
**Goal**: Extract business logic from handlers

#### Tasks:
1. ✅ Create service layer
   - Create `internal/service/niv_service.go`
   - Move business logic from handlers
   - Create `internal/service/explain_service.go`
   - Extract OpenAI API calls from `server/niv_server.go`

2. ✅ Refactor handlers
   - Update handlers to only handle HTTP concerns
   - Handlers call services, not repositories
   - Proper error handling and HTTP status codes

3. ✅ Add service tests
   - Unit tests for services
   - Mock repositories in tests

**Deliverables**:
- ✅ Service layer implemented
- ✅ Business logic separated from HTTP layer
- ✅ OpenAI integration moved to service
- ✅ Service tests passing

---

### Phase 4: API Layer Refactoring (Week 2-3)
**Goal**: Organize HTTP handlers and middleware

#### Tasks:
1. ✅ Move handlers to `internal/api/handlers/`
   - Split `server/niv_server.go` into focused handlers
   - `health.go`: Health check handlers
   - `niv_handler.go`: NIV-related handlers

2. ✅ Extract middleware
   - Move CORS to `internal/api/middleware/cors.go`
   - Create `internal/api/middleware/logger.go`
   - Create `internal/api/middleware/recovery.go`

3. ✅ Create routes package
   - Move route registration to `internal/api/routes/routes.go`
   - Clean separation of route definitions

4. ✅ Refactor server setup
   - Create `internal/api/server.go` (or keep in `cmd/api/main.go`)
   - Initialize Echo server with middleware and routes

**Deliverables**:
- ✅ Handlers organized by feature
- ✅ Middleware extracted and reusable
- ✅ Routes centralized
- ✅ Clean server initialization

---

### Phase 5: DTOs & Models Organization (Week 3)
**Goal**: Organize data structures

#### Tasks:
1. ✅ Organize DTOs
   - Move to `dto/request/` and `dto/response/`
   - Separate request/response DTOs
   - Add validation tags where needed

2. ✅ Review models
   - Ensure models are pure domain entities
   - Add any missing models

3. ✅ Update imports
   - Update all imports to new DTO locations

**Deliverables**:
- ✅ DTOs organized by type
- ✅ Models clean and focused
- ✅ All imports updated

---

### Phase 6: Testing & Scripts (Week 3-4)
**Goal**: Improve testing and automation

#### Tasks:
1. ✅ Organize tests
   - Keep unit tests with code (`*_test.go`)
   - Move integration tests to `test/integration/`
   - Add test fixtures if needed

2. ✅ Create scripts
   - `scripts/build.sh`: Build script
   - `scripts/test.sh`: Test runner
   - `scripts/migrate.sh`: Migration helper
   - `scripts/docker-build.sh`: Docker build

3. ✅ Create Makefile
   - Common commands (build, test, run, migrate)
   - Docker commands

**Deliverables**:
- ✅ Tests organized
- ✅ Scripts created
- ✅ Makefile with common commands

---

### Phase 7: CI/CD & Documentation (Week 4)
**Goal**: Automation and documentation

#### Tasks:
1. ✅ Create CI/CD workflows
   - `.github/workflows/ci.yml`: Run tests on PR
   - `.github/workflows/docker.yml`: Build and push Docker image

2. ✅ Update Dockerfile
   - Update build path to `./cmd/api`
   - Optimize for new structure

3. ✅ Create documentation
   - Update `README.md` with new structure
   - Create `docs/API.md` (API documentation)
   - Create `docs/DEPLOYMENT.md` (deployment guide)

4. ✅ Create `.dockerignore`
   - Exclude unnecessary files from Docker builds

**Deliverables**:
- ✅ CI/CD pipelines working
- ✅ Dockerfile updated
- ✅ Documentation complete
- ✅ `.dockerignore` created

---

### Phase 8: Cleanup & Validation (Week 4)
**Goal**: Final cleanup and verification

#### Tasks:
1. ✅ Remove old files
   - Delete root `main.go` (after migration)
   - Clean up any unused files

2. ✅ Verify all imports
   - Ensure all imports are correct
   - Run `go mod tidy`

3. ✅ Run full test suite
   - Unit tests
   - Integration tests
   - Manual testing

4. ✅ Update `.gitignore`
   - Add any new patterns needed

5. ✅ Code review
   - Review structure
   - Check for any remaining issues

**Deliverables**:
- ✅ Old files removed
- ✅ All tests passing
- ✅ Code reviewed
- ✅ Project ready for production

---

## 📊 Migration Checklist

### Phase 1: Foundation
- [ ] Create directory structure
- [ ] Create `internal/config/` package
- [ ] Move config logic from database client
- [ ] Create `cmd/api/main.go`
- [ ] Update imports
- [ ] Project compiles

### Phase 2: Database & Repository
- [ ] Refactor `internal/database/client.go`
- [ ] Create repository interfaces
- [ ] Implement NIV repository
- [ ] Move queries to repository
- [ ] Update repository tests
- [ ] All repository tests passing

### Phase 3: Service Layer
- [ ] Create NIV service
- [ ] Create explain service
- [ ] Move OpenAI logic to service
- [ ] Update handlers to use services
- [ ] Add service tests
- [ ] All service tests passing

### Phase 4: API Layer
- [ ] Move handlers to `internal/api/handlers/`
- [ ] Extract middleware
- [ ] Create routes package
- [ ] Refactor server setup
- [ ] All handlers working

### Phase 5: DTOs & Models
- [ ] Organize DTOs into request/response
- [ ] Review models
- [ ] Update all imports
- [ ] All DTOs organized

### Phase 6: Testing & Scripts
- [ ] Organize tests
- [ ] Create build scripts
- [ ] Create Makefile
- [ ] Scripts tested

### Phase 7: CI/CD & Docs
- [ ] Create CI workflow
- [ ] Create Docker workflow
- [ ] Update Dockerfile
- [ ] Create documentation
- [ ] Create `.dockerignore`

### Phase 8: Cleanup
- [ ] Remove old files
- [ ] Verify imports
- [ ] Run full test suite
- [ ] Code review
- [ ] Ready for production

---

## 🎯 Benefits of New Structure

### Maintainability
- ✅ Clear separation of concerns
- ✅ Easy to locate code
- ✅ Consistent patterns

### Testability
- ✅ Each layer can be tested independently
- ✅ Easy to mock dependencies
- ✅ Integration tests separate from unit tests

### Scalability
- ✅ Easy to add new features
- ✅ Multiple entry points possible
- ✅ Clear boundaries for microservices

### Developer Experience
- ✅ Follows Go conventions
- ✅ Easy onboarding for new developers
- ✅ Clear project layout

---

## 📚 References

- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Go Project Structure Best Practices](https://github.com/golang/go/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)

---

## 🔄 Next Steps

1. Review this document with the team
2. Start with Phase 1 (Foundation)
3. Migrate incrementally, one phase at a time
4. Test thoroughly after each phase
5. Update this document as needed

---

**Last Updated**: 2025-01-27
**Version**: 1.0

