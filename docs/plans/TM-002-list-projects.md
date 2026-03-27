# Implementation Plan: TM-002 - List Projects

## Current State
- Backend has stub `List` method (returns nil)
- Handler has pagination but page param parsing bug (line 113 doesn't assign value)
- Default page size is 10 (should be 20)
- No sorting implemented
- No search implemented
- No task count field
- Frontend has basic list view but no search/pagination UI

## Implementation Checklist

### Phase 1: Backend - Domain Layer
- [ ] Add `task_count` field to `Project` struct in `backend/internal/domain/project/project.go`
- [ ] Update `ProjectListFilter` to use default page_size of 20

### Phase 2: Backend - Storage Layer
- [ ] Implement `List` method in `backend/internal/storage/postgres/project_repository.go`:
  - Query projects with `WHERE deleted_at IS NULL`
  - Add `ORDER BY created_at DESC` for newest first
  - Add `ILIKE '%search%'` for name search filter
  - Add pagination with `LIMIT` and `OFFSET`
  - Return total count for pagination
  - Return `task_count` as 0 (placeholder for tasks feature)

### Phase 3: Backend - Handler Layer
- [ ] Fix page parameter parsing bug in `project_handler.go:113`
- [ ] Update default page_size to 20
- [ ] Parse page_size from query param

### Phase 4: Backend - Tests
- [ ] Add unit tests for `List` usecase method
- [ ] Add unit tests for `ListProjects` handler
- [ ] Add integration tests for repository `List` method

### Phase 5: Frontend - Types
- [ ] Add `task_count` to `Project` interface in `frontend/lib/types/project.ts`

### Phase 6: Frontend - UI Components
- [ ] Add search input component to projects page
- [ ] Add pagination component
- [ ] Display task count on project cards
- [ ] Wire up search/filter to API calls

### Phase 7: Frontend - Tests
- [ ] Add tests for search functionality
- [ ] Add tests for pagination

## Files to Modify
| File | Changes |
|------|---------|
| `backend/internal/domain/project/project.go` | Add task_count, update filter defaults |
| `backend/internal/storage/postgres/project_repository.go` | Implement List method |
| `backend/internal/usecase/project/project_usecase.go` | Add List tests |
| `backend/internal/handler/project/project_handler.go` | Fix page bug, update defaults |
| `frontend/lib/types/project.ts` | Add task_count field |
| `frontend/app/projects/page.tsx` | Add search & pagination UI |
| `frontend/lib/api/projects.ts` | Already has filter support ✓ |

## API Contract
```
GET /api/v1/projects?page=1&page_size=20&search=keyword

Response 200:
{
  "data": [
    {
      "id": "uuid",
      "name": "Project Name",
      "description": "Description",
      "status": "active",
      "task_count": 0,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "total": 42
}
```