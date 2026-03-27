# Implementation Plan: TM-001 Create Project

## Overview
Implement the "Create Project" feature following TDD (Test-Driven Development) with Red-Green-Refactor cycles. This spans both backend (Go/Gin) and frontend (Next.js 14).

---

## Phase 0: Infrastructure Setup

### 0.1 Create Docker Compose for PostgreSQL
- Create `docker-compose.yml` at project root
- Services: PostgreSQL 16
- Configure volumes for data persistence
- Configure environment variables for database connection
- Add health checks for PostgreSQL readiness

### 0.2 Database Configuration
- Update `.env.example` with database URL template
- Ensure backend config reads from environment

---

## Phase 1: Backend Storage Layer (Database + Repository)

### 1.1 Create Database Migration
- Create `backend/migrations/001_create_projects_table.sql`
- Fields: `id` (UUID), `name`, `description`, `status`, `created_at`, `updated_at`, `deleted_at`
- Add unique constraint on `name`

### 1.2 Implement Project Repository (Red → Green → Refactor)
**Test First (Red):**
- Create `backend/internal/storage/postgres/project_repository_test.go`
- Test cases:
  - Success: Create project with valid input
  - Validation: Duplicate name returns `ErrProjectNameExists`
  - Infrastructure: Database connection errors

**Implement (Green):**
- Create `backend/internal/storage/postgres/project_repository.go`
- Implement `Create(ctx, *Project) error`
- Implement `GetByName(ctx, name) (*Project, error)` for uniqueness check
- Use UUID generation for ID

---

## Phase 2: Backend Usecase Layer

### 2.1 Implement Project Usecase (Red → Green → Refactor)
**Test First (Red):**
- Create `backend/internal/usecase/project_usecase_test.go`
- Test cases:
  - Success: Create project returns project with generated ID
  - Success: Status defaults to "active"
  - Success: Timestamps are recorded
  - Validation: Name required (max 100 chars)
  - Validation: Description max 500 chars
  - Business: Duplicate name returns `ErrProjectNameExists`
  - Context: Cancellation handling

**Implement (Green):**
- Create `backend/internal/usecase/project_usecase.go`
- Inject Repository interface
- Validate input using validator
- Check unique name before create
- Set default status and timestamps

---

## Phase 3: Backend Handler Layer

### 3.1 Implement Create Project Handler (Red → Green → Refactor)
**Test First (Red):**
- Create `backend/internal/handler/project_handler_test.go`
- Test cases:
  - Success: POST /api/v1/projects returns 201 with project
  - Validation: Missing name returns 400
  - Validation: Name > 100 chars returns 400
  - Business: Duplicate name returns 409
  - Infrastructure: Internal error returns 500

**Implement (Green):**
- Update `backend/internal/handler/router.go`
- Create `backend/internal/handler/project_handler.go`
- Inject Usecase interface
- Parse JSON request body
- Call usecase.CreateProject
- Return appropriate HTTP status codes

---

## Phase 4: Backend Integration

### 4.1 Wire Dependencies
- Update `backend/cmd/server/main.go`
- Initialize DB connection
- Initialize repository, usecase, handler
- Start Gin server with routes

### 4.2 Integration Test
- Create `backend/integration/project_test.go` (optional)
- Test full flow with testcontainers

---

## Phase 5: Frontend API Layer

### 5.1 Create Project API Hook (Red → Green → Refactor)
**Test First (Red):**
- Create `frontend/lib/api/projects.test.ts`
- Test API client mutation

**Implement (Green):**
- Update `frontend/lib/api.ts` with project endpoints
- Create `frontend/lib/hooks/useProjects.ts`
- Implement `useCreateProject` mutation hook with React Query

---

## Phase 6: Frontend UI Components

### 6.1 Create Project Modal (Red → Green → Refactor)
**Test First (Red):**
- Create `frontend/components/projects/CreateProjectModal.test.tsx`
- Test cases:
  - Renders form with name and description fields
  - Shows validation errors for invalid input
  - Calls mutation on submit
  - Shows loading state during submission
  - Closes modal on success

**Implement (Green):**
- Create `frontend/components/ui/Modal.tsx` (reusable)
- Create `frontend/components/ui/Input.tsx` (reusable)
- Create `frontend/components/ui/Button.tsx` (reusable)
- Create `frontend/components/projects/CreateProjectModal.tsx`
- Implement form with React Hook Form
- Add validation (name required, max chars)
- Handle loading/error states

### 6.2 Create Projects Page
- Create `frontend/app/projects/page.tsx`
- Add "Create Project" button
- Integrate CreateProjectModal

---

## Phase 7: End-to-End Verification

### 7.1 Manual Testing Checklist
- [ ] Backend: POST /api/v1/projects creates project
- [ ] Backend: Returns 201 with project details
- [ ] Backend: Validates unique name
- [ ] Frontend: Modal opens/closes correctly
- [ ] Frontend: Form validation works
- [ ] Frontend: Project created on submit

---

## Files to Create/Modify

### Infrastructure (New Files)
```
docker-compose.yml
.env.example
```

### Backend (New Files)
```
backend/
├── migrations/
│   └── 001_create_projects_table.sql
├── internal/
│   ├── storage/postgres/
│   │   ├── project_repository.go
│   │   └── project_repository_test.go
│   ├── usecase/
│   │   ├── project_usecase.go
│   │   └── project_usecase_test.go
│   └── handler/
│       ├── project_handler.go
│       └── project_handler_test.go
```

### Backend (Modified Files)
```
backend/
├── cmd/server/main.go          # Wire dependencies
├── internal/handler/router.go  # Update createProject handler
```

### Frontend (New Files)
```
frontend/
├── components/
│   ├── ui/
│   │   ├── Modal.tsx
│   │   ├── Input.tsx
│   │   └── Button.tsx
│   └── projects/
│       └── CreateProjectModal.tsx
├── lib/
│   └── hooks/
│       └── useProjects.ts
└── app/
    └── projects/
        └── page.tsx
```

### Frontend (Modified Files)
```
frontend/
└── lib/api.ts  # Add project API functions
```

---

## Estimated Effort
- Infrastructure (Docker): ~30 min
- Backend Storage: ~1.5 hours
- Backend Usecase: ~1 hour
- Backend Handler: ~1 hour
- Frontend API: ~30 min
- Frontend UI: ~1.5 hours
- Integration: ~30 min
- **Total: ~6.5 hours**

---

## TDD Workflow
For each component:
1. **RED**: Write failing test first
2. **GREEN**: Write minimal code to pass
3. **REFACTOR**: Clean up while keeping tests green