# Implementation Plan: TM-003 - Get Project Details

## Current State

**Backend (Already Implemented):**
- `GET /api/v1/projects/:id` endpoint exists in `project_handler.go:77-105`
- `GetProject` usecase method exists in `project_usecase.go:74-83`
- `GetByID` repository method exists in `project_repository.go:65-95`
- Returns 404 for not found, 200 with project details

**Frontend (Missing):**
- No project detail page
- `useProject` hook exists in `lib/api/projects.ts`
- Project cards have "View" button but no navigation

**Gaps:**
- `task_count` not included in single project response
- `completion_percentage` not in Project struct
- No frontend detail page

---

## Phase 1: Backend - Add Task Statistics

### 1.1 Update Project Entity
- Add `CompletionPercentage` field to `Project` struct in `backend/internal/domain/project/project.go`

### 1.2 Update Repository GetByID
- Modify `GetByID` in `backend/internal/storage/postgres/project_repository.go`
- Add subquery to count tasks for the project
- Calculate completion percentage (done_tasks / total_tasks * 100)
- Return `task_count` and `completion_percentage` in response

### 1.3 Add Tests
- Unit test for `GetProject` usecase
- Unit test for `GetProject` handler
- Integration test for `GetByID` repository

---

## Phase 2: Frontend - Types & API

### 2.1 Update Types
- Add `completion_percentage` to `Project` interface in `frontend/lib/types/project.ts`

### 2.2 Verify API Hook
- `useProject(id)` hook already exists in `frontend/lib/api/projects.ts` ✓

---

## Phase 3: Frontend - Project Detail Page

### 3.1 Create Detail Page
- Create `frontend/app/projects/[id]/page.tsx`
- Fetch project using `useProject(id)` hook
- Display project name, description, status badge
- Display created/updated timestamps (formatted)
- Display task count and completion percentage with progress bar
- Handle loading state with skeleton
- Handle error/not found state with EmptyState
- Add "Back to Projects" navigation

### 3.2 Update Projects List Page
- Make project cards clickable (link to detail page)
- Add navigation on "View" button click

### 3.3 Add Tests
- Test project detail page rendering
- Test loading/error states
- Test navigation from list to detail

---

## API Contract

```
GET /api/v1/projects/:id

Response 200:
{
  "id": "uuid",
  "name": "Project Name",
  "description": "Description",
  "status": "active",
  "task_count": 5,
  "completion_percentage": 60.0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}

Response 404:
{
  "error": "not_found",
  "message": "project not found"
}
```

---

## Files to Modify

| File | Changes |
|------|---------|
| `backend/internal/domain/project/project.go` | Add `CompletionPercentage` field |
| `backend/internal/storage/postgres/project_repository.go` | Add task count query in GetByID |
| `frontend/lib/types/project.ts` | Add `completion_percentage` field |
| `frontend/app/projects/[id]/page.tsx` | **Create** - Project detail page |
| `frontend/app/projects/page.tsx` | Add navigation to detail page |

---

## Estimated Effort
- Backend (Repository + Tests): ~1 hour
- Frontend (Page + Navigation): ~1.5 hours
- **Total: ~2.5 hours**