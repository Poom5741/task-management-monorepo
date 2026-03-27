# Backend TDG.md

## Test Configuration

### Build Command
```bash
make build
```

### Test Command (Full Suite)
```bash
make test
```

### Single Test Command
```bash
go test -run TestHandler_MethodName ./internal/handler/project
```

### Coverage Command
```bash
make test-coverage
```

### Test File Patterns

#### Handler Tests
- Location: `internal/handler/{domain}/{domain}_handler_test.go`
- Pattern: `Test{Handler}_{Method}`

#### Usecase Tests
- Location: `internal/usecase/{domain}/{domain}_usecase_test.go`
- Pattern: `Test{Usecase}_{Method}`

#### Storage Tests
- Location: `internal/storage/postgres/{domain}_repository_test.go`
- Pattern: `Test{Repository}_{Method}`

### Coverage Threshold
- Minimum: 70%
- Target: 85%

### Testcontainers Configuration
```go
func setupTestContainers(t *testing.T) (*postgres.DB, func()) {
    ctx := context.Background()
    
    pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "postgres:15",
            Env: map[string]string{
                "POSTGRES_DB":      "testdb",
                "POSTGRES_USER":    "test",
                "POSTGRES_PASSWORD": "test",
            },
            ExposedPorts: []string{"5432/tcp"},
        },
    })
    // ...
}
```

### Mock Generation
Use Function Field Mock Pattern (no code generation needed):

```go
type mockProjectRepository struct {
    CreateFunc    func(ctx context.Context, p *project.Project) error
    GetByIDFunc   func(ctx context.Context, id string) (*project.Project, error)
    ListFunc      func(ctx context.Context, filter *project.ProjectListFilter) ([]*project.Project, int, error)
    UpdateFunc    func(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error)
    DeleteFunc    func(ctx context.Context, id string) error
}
```

### Pre-commit Hooks
Run before every commit:
1. `make fmt` - Format code
2. `make lint` - Run linter
3. `make test` - Run tests