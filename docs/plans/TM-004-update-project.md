# Implementation Plan: TM-004 - Update Project

## Current State Analysis

**Backend (Partially Implemented):**
- ✅ `UpdateProjectInput` struct exists in `backend/internal/domain/project/project.go:30-34`
- ✅ `Update` interface method defined in `backend/internal/domain/project/interfaces.go:12`
- ✅ `UpdateProject` handler skeleton exists in `backend/internal/handler/project/project_handler.go:144-181`
- ✅ `UpdateProject` usecase skeleton exists in `backend/internal/usecase/project/project_usecase.go:89-91`
- ❌ `Update` repository method returns `nil, nil` (not implemented) in `backend/internal/storage/postgres/project_repository.go:211-213`

**Frontend (Ready for Implementation):**
- ✅ `UpdateProjectInput` type exists in `frontend/lib/types/project.ts:17-21`
- ✅ `useUpdateProject` hook exists in `frontend/lib/api/projects.ts:62-73`
- ❌ No EditProjectModal component
- ❌ No edit button on project detail page

---

## Phase 1: Backend - Repository Layer (TDD)

### 1.1 Write Repository Tests (Red)
Create tests in `backend/internal/storage/postgres/project_repository_test.go`:

**Test Categories:**
- **Success:** Update project name only
- **Success:** Update description only
- **Success:** Update status to archived
- **Success:** Update multiple fields at once
- **Validation:** Update name to existing name returns `ErrProjectNameExists`
- **Infrastructure:** Project not found returns `ErrProjectNotFound`
- **Context:** Database connection errors

### 1.2 Implement Update Method (Green)
Implement `Update` in `backend/internal/storage/postgres/project_repository.go`:

```go
func (r *ProjectRepository) Update(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
    // 1. Check if project exists
    // 2. Build dynamic UPDATE query based on non-nil fields
    // 3. Always update updated_at timestamp
    // 4. Execute query with conflict handling for unique name
    // 5. Return updated project with task statistics
}
```

**Implementation Details:**
- Use `COALESCE` pattern for partial updates
- Add `WHERE id = $id AND deleted_at IS NULL` condition
- Handle unique constraint violation (23505) → return `ErrProjectNameExists`
- Return `ErrProjectNotFound` if no rows affected
- Reuse `GetByID` query pattern to return updated project with task statistics

---

## Phase 2: Backend - Usecase Layer (TDD)

### 2.1 Write Usecase Tests (Red)
Create tests in `backend/internal/usecase/project/project_usecase_test.go`:

**Test Categories:**
- **Success:** Update project returns updated project
- **Success:** Update only name (other fields remain unchanged)
- **Success:** Update only description
- **Success:** Update only status
- **Validation:** Empty ID returns validation error
- **Validation:** Name > 100 chars returns validation error
- **Validation:** Description > 500 chars returns validation error
- **Validation:** Invalid status returns validation error
- **Business:** Update name to existing name returns `ErrProjectNameExists`
- **Business:** Project not found returns `ErrProjectNotFound`
- **Context:** Cancellation handling

### 2.2 Implement UpdateProject Method (Green)
Implement `UpdateProject` in `backend/internal/usecase/project/project_usecase.go`:

```go
func (uc *ProjectUsecase) UpdateProject(ctx context.Context, id string, input *project.UpdateProjectInput) (*project.Project, error) {
    // 1. Validate input (nil check, field validations)
    // 2. Validate status if provided
    // 3. Call repository.Update
    // 4. Return updated project
}
```

**Implementation Details:**
- Validate ID is not empty
- Validate name length if provided
- Validate description length if provided
- Validate status is valid enum if provided
- Check name uniqueness if name is being updated
- Delegate to repository

---

## Phase 3: Backend - Handler Layer (TDD)

### 3.1 Write Handler Tests (Red)
Add tests to `backend/internal/handler/project/project_handler_test.go`:

**Test Categories:**
- **Success:** PATCH /api/v1/projects/:id returns 200 with updated project
- **Success:** Update name only
- **Success:** Update status to archived
- **Validation:** Missing ID returns 400
- **Validation:** Invalid JSON returns 400
- **Validation:** Name > 100 chars returns 400
- **Validation:** Invalid status returns 400
- **Business:** Duplicate name returns 409 Conflict
- **Business:** Project not found returns 404
- **Infrastructure:** Internal error returns 500

### 3.2 Verify Handler Implementation (Green)
The handler at `backend/internal/handler/project/project_handler.go:144-181` is mostly complete but needs:

**Missing Error Handling:**
- Add conflict error handling for `ErrProjectNameExists` → return 409
- Add validation error handling → return 400

```go
// Add to error handling block:
if errors.Is(err, project.ErrProjectNameExists) {
    c.JSON(http.StatusConflict, ErrorResponse{
        Error:   "conflict",
        Message: "project name already exists",
    })
    return
}

var validationErr *project.ValidationError
if errors.As(err, &validationErr) {
    c.JSON(http.StatusBadRequest, ErrorResponse{
        Error:   "validation_error",
        Message: validationErr.Message,
    })
    return
}
```

---

## Phase 4: Frontend - Edit Project Modal

### 4.1 Create EditProjectModal Component
Create `frontend/components/projects/EditProjectModal.tsx`:

**Component Structure:**
```typescript
interface EditProjectModalProps {
  isOpen: boolean
  onClose: () => void
  project: Project  // Pre-populate form with existing data
}
```

**Features:**
- Pre-populate form with current project data
- Form fields: name, description, status (dropdown)
- Client-side validation (same as CreateProjectModal)
- Handle 409 conflict error (duplicate name)
- Handle 404 not found error
- Show loading state during mutation
- Call `useUpdateProject` hook on submit
- Invalidate queries on success (already handled by hook)
- Close modal on success

**Validation Rules:**
- Name: required, max 100 chars
- Description: optional, max 500 chars
- Status: must be 'active' or 'archived'

### 4.2 Create Modal Tests
Create `frontend/components/projects/EditProjectModal.test.tsx`:

**Test Cases:**
- Renders with pre-populated data
- Shows validation errors for invalid input
- Calls update mutation on submit
- Shows loading state during submission
- Handles 409 conflict error (duplicate name)
- Handles 404 not found error
- Closes modal on success
- Status dropdown shows both options

---

## Phase 5: Frontend - Integration

### 5.1 Update Project Detail Page
Modify `frontend/app/projects/[id]/page.tsx`:

**Changes:**
- Add "Edit" button next to status badge (lines 95-100)
- Add state for modal open/close
- Import and render `EditProjectModal`
- Pass current project to modal

**Button Placement:**
```typescript
<div className="flex items-start justify-between">
  <div>
    <h1>{project.name}</h1>
    <p>Created/Updated dates</p>
  </div>
  <div className="flex items-center gap-3">
    <Button
      variant="secondary"
      onClick={() => setIsEditModalOpen(true)}
      leftIcon={<EditIcon />}
    >
      Edit
    </Button>
    <Badge variant={...}>{project.status}</Badge>
  </div>
</div>
```

### 5.2 Update Detail Page Tests
Modify `frontend/app/projects/[id]/page.test.tsx`:

**Add Test Cases:**
- Edit button renders and is clickable
- Modal opens when Edit button clicked
- Modal closes on cancel
- Project data refreshes after successful update

---

## Phase 6: Integration Testing

### 6.1 Backend Integration Tests
Create `backend/integration/project_update_test.go` (optional):

**Test Scenarios:**
- Full flow: Create → Update → Verify
- Update with duplicate name
- Update non-existent project

### 6.2 Manual E2E Testing Checklist
- [ ] Backend: PATCH /api/v1/projects/:id updates project
- [ ] Backend: Returns 200 with updated project
- [ ] Backend: Returns 404 for non-existent project
- [ ] Backend: Returns 409 for duplicate name
- [ ] Backend: Updates `updated_at` timestamp
- [ ] Backend: Validates unique name constraint
- [ ] Frontend: Edit modal opens with pre-populated data
- [ ] Frontend: Form validation works correctly
- [ ] Frontend: Can update each field independently
- [ ] Frontend: Shows error for duplicate name
- [ ] Frontend: Project details refresh after update
- [ ] Frontend: Status dropdown works correctly

---

## API Contract

```
PATCH /api/v1/projects/:id
Content-Type: application/json

Request Body (all fields optional):
{
  "name": "Updated Project Name",
  "description": "Updated description",
  "status": "archived"
}

Response 200:
{
  "id": "uuid",
  "name": "Updated Project Name",
  "description": "Updated description",
  "status": "archived",
  "task_count": 5,
  "completion_percentage": 60.0,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-02T12:30:00Z"  // Updated timestamp
}

Response 400:
{
  "error": "validation_error",
  "message": "name must be at most 100 characters"
}

Response 404:
{
  "error": "not_found",
  "message": "project not found"
}

Response 409:
{
  "error": "conflict",
  "message": "project name already exists"
}
```

---

## Files to Create/Modify

### Backend Files

| File | Status | Changes |
|------|--------|---------|
| `backend/internal/storage/postgres/project_repository.go` | Modify | Implement `Update` method (lines 211-213) |
| `backend/internal/storage/postgres/project_repository_test.go` | Modify | Add Update test cases |
| `backend/internal/usecase/project/project_usecase.go` | Modify | Implement `UpdateProject` validation (lines 89-91) |
| `backend/internal/usecase/project/project_usecase_test.go` | Modify | Add UpdateProject test cases |
| `backend/internal/handler/project/project_handler.go` | Modify | Add conflict & validation error handling (lines 164-178) |
| `backend/internal/handler/project/project_handler_test.go` | Modify | Add UpdateProject handler tests |

### Frontend Files

| File | Status | Changes |
|------|--------|---------|
| `frontend/components/projects/EditProjectModal.tsx` | **Create** | New modal component |
| `frontend/components/projects/EditProjectModal.test.tsx` | **Create** | Modal tests |
| `frontend/app/projects/[id]/page.tsx` | Modify | Add Edit button and modal integration |
| `frontend/app/projects/[id]/page.test.tsx` | Modify | Add edit functionality tests |

---

## Estimated Effort

- **Phase 1:** Backend Repository (Tests + Implementation) - 1.5 hours
- **Phase 2:** Backend Usecase (Tests + Implementation) - 1 hour
- **Phase 3:** Backend Handler (Tests + Error Handling) - 1 hour
- **Phase 4:** Frontend Modal (Component + Tests) - 2 hours
- **Phase 5:** Frontend Integration - 0.5 hours
- **Phase 6:** Integration Testing - 0.5 hours
- **Total:** ~6.5 hours

---

## TDD Workflow

For each component:
1. **RED:** Write failing test first
2. **GREEN:** Write minimal code to pass
3. **REFACTOR:** Clean up while keeping tests green

---

## Dependencies & Risks

**Dependencies:**
- ✅ Database migration already exists (001_create_projects_table.sql)
- ✅ API types already defined
- ✅ React Query hooks already implemented

**Risks:**
- Name uniqueness check on update needs to exclude current project
- Status enum validation must match backend constants
- Need to handle race condition for name updates

**Mitigation:**
- Use `WHERE name = $name AND id != $id` for uniqueness check
- Use backend Status constants in validation
- Use database unique constraint for consistency