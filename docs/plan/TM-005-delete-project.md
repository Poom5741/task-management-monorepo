# TM-005: Delete Project Implementation Plan

## Overview
Implement soft delete for projects with cascade deletion of associated tasks. Includes backend API implementation and frontend confirmation dialog.

## Current State
- **Backend**: Handler implemented, usecase stub, repository stub
- **Frontend**: `useDeleteProject` hook exists, no confirmation modal, no delete button

---

## Phase 1: Backend Implementation

### 1.1 Implement Repository Delete Method
**File**: `backend/internal/storage/postgres/project_repository.go`

```go
func (r *ProjectRepository) Delete(ctx context.Context, id string) error {
    // 1. Verify project exists
    _, err := r.GetByID(ctx, id)
    if err != nil {
        return err
    }
    
    // 2. Soft delete tasks (cascade)
    taskQuery := `UPDATE tasks SET deleted_at = NOW() WHERE project_id = $1 AND deleted_at IS NULL`
    r.db.ExecContext(ctx, taskQuery, id)
    
    // 3. Soft delete project
    query := `UPDATE projects SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
    result, err := r.db.ExecContext(ctx, query, id)
    if err != nil {
        return err
    }
    
    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        return project.ErrProjectNotFound
    }
    
    return nil
}
```

### 1.2 Update Usecase (Add Validation)
**File**: `backend/internal/usecase/project/project_usecase.go`

Add ID validation before calling repository:
```go
func (uc *ProjectUsecase) DeleteProject(ctx context.Context, id string) error {
    if id == "" {
        return &project.ValidationError{
            Field:   "id",
            Message: "id is required",
        }
    }
    return uc.repo.Delete(ctx, id)
}
```

### 1.3 Create Tasks Table Migration (if needed)
**File**: `backend/migrations/002_create_tasks_table.sql`

Since the repository queries tasks, ensure the table exists with proper foreign key and `deleted_at` column.

---

## Phase 2: Frontend Implementation

### 2.1 Create DeleteProjectModal Component
**File**: `frontend/components/projects/DeleteProjectModal.tsx`

- Warning icon with danger styling
- Confirmation message: "Are you sure you want to delete this project?"
- Warning text: "This will also delete all tasks associated with this project. This action cannot be undone."
- Cancel (ghost) + Delete (danger) buttons
- Loading state during deletion
- Error handling

### 2.2 Add Delete Button to Project Detail Page
**File**: `frontend/app/projects/[id]/page.tsx`

- Add Delete button next to Edit button (danger variant)
- Wire up DeleteProjectModal
- Redirect to `/projects` on successful deletion
- Handle 404 error (project already deleted)

---

## Phase 3: Testing

### 3.1 Backend Unit Tests
**File**: `backend/internal/handler/project/project_handler_test.go`

Test cases:
1. **Success** - Delete returns 204
2. **Not Found** - Non-existent project returns 404
3. **Validation** - Empty ID returns 400

### 3.2 Frontend Component Tests
**File**: `frontend/components/projects/DeleteProjectModal.test.tsx`

Test cases:
1. Renders confirmation dialog
2. Cancel closes modal without deleting
3. Delete calls mutation and closes on success
4. Shows error message on failure

---

## Files to Modify/Create

| Action | File |
|--------|------|
| Modify | `backend/internal/storage/postgres/project_repository.go` |
| Modify | `backend/internal/usecase/project/project_usecase.go` |
| Create | `backend/migrations/002_create_tasks_table.sql` (if needed) |
| Create | `frontend/components/projects/DeleteProjectModal.tsx` |
| Modify | `frontend/app/projects/[id]/page.tsx` |
| Create | `frontend/components/projects/DeleteProjectModal.test.tsx` |

---

## Acceptance Criteria Mapping

| Criteria | Implementation |
|----------|---------------|
| System confirms deletion | DeleteProjectModal with confirmation |
| All tasks deleted | Repository cascade soft-delete |
| Project soft-deleted | Repository sets `deleted_at` |
| Return 404 if not found | Handler returns 404 (already implemented) |
| API returns 204 | Handler returns 204 (already implemented) |