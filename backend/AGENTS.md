# Backend AGENTS.md

## Project Overview

Task Management System Backend - Go API with Hexagonal Architecture

## Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters) pattern:

```
cmd/server/          # Application entry point
internal/
├── domain/          # Core domain layer (entities, business rules)
│   ├── project/
│   ├── task/
│   ├── label/
│   └── dependency/
├── handler/         # HTTP delivery layer (adapters)
├── usecase/         # Application layer (business logic orchestration)
└── storage/         # Infrastructure layer (data access adapters)
    ├── postgres/
    └── redis/
pkg/                 # Shared packages
```

## Domain Structure

Each domain follows this structure:
- `{domain}.go` - Entity definitions and value objects
- `interfaces.go` - Repository and Usecase interfaces
- `errors.go` - Domain-specific errors

## Testing Patterns

### Handler Tests
Use **Function Field Mock Pattern**:
```go
type mockProjectUsecase struct {
    CreateProjectFunc func(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error)
}

func (m *mockProjectUsecase) CreateProject(ctx context.Context, input *project.CreateProjectInput) (*project.Project, error) {
    return m.CreateProjectFunc(ctx, input)
}
```

### Test Categories
Every test must cover 5 categories:
1. **Success** - Happy path
2. **Context** - Cancellation, timeout
3. **Validation** - Input validation errors
4. **Business** - Business rule violations
5. **Infrastructure** - Database, network errors

### Testcontainers
Use testcontainers-go for integration tests with real PostgreSQL/Redis.

## Code Conventions

### Handler Layer
- Gin framework for HTTP routing
- JSON request/response
- HTTP status codes:
  - 200: Success (GET, PATCH)
  - 201: Created (POST)
  - 204: No Content (DELETE)
  - 400: Validation error
  - 404: Not found
  - 500: Internal error

### Usecase Layer
- Pure business logic (no HTTP dependencies)
- Return domain errors
- Use context for cancellation

### Storage Layer
- PostgreSQL for persistence
- Redis for caching
- Implement Repository interface

## Anti-Patterns (Avoid)

- ❌ Import cycles
- ❌ Business logic in handlers
- ❌ HTTP concerns in usecase/storage
- ❌ Direct database queries in handlers
- ❌ Hardcoded values (use config)

## Commands

```bash
# Run server
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build binary
make build

# Format code
make fmt

# Run linter
make lint
```

## Error Handling

Define domain-specific errors in each domain's `errors.go`:
```go
var (
    ErrProjectNotFound   = errors.New("project not found")
    ErrProjectNameExists = errors.New("project name already exists")
)
```

## Validation

Use struct tags for validation:
```go
type CreateProjectInput struct {
    Name        string `json:"name" validate:"required,max=100"`
    Description string `json:"description" validate:"max=500"`
}
```