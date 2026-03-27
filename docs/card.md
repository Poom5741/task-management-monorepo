# Task Management System - Feature Cards

## Epic 1: Project Management

### TM-001: Create Project Done
**As a** user  
**I want to** create a new project  
**So that** I can organize related tasks together  

**Acceptance Criteria:**
- [ ] User can input project name (required, max 100 chars)
- [ ] User can input project description (optional, max 500 chars)
- [ ] System validates unique project name per user
- [ ] System generates unique project ID
- [ ] Project is created with "active" status by default
- [ ] Created timestamp is recorded
- [ ] API returns 201 with project details

**Technical Notes:**
- Backend: POST /api/v1/projects
- Frontend: Create project modal form

---

### TM-002: List Projects Done 
**As a** user  
**I want to** view all my projects  
**So that** I can see and select projects to work on  

**Acceptance Criteria:**
- [ ] Display all projects with name, description, task count
- [ ] Show projects sorted by created date (newest first)
- [ ] Support pagination (20 items per page)
- [ ] Support search/filter by project name
- [ ] API returns 200 with projects array

**Technical Notes:**
- Backend: GET /api/v1/projects
- Frontend: Project list page with search

---

### TM-003: Get Project Details
**As a** user  
**I want to** view a specific project  
**So that** I can see all its details  

**Acceptance Criteria:**
- [ ] Display project name, description, status
- [ ] Display created and updated timestamps
- [ ] Display task count and completion percentage
- [ ] Return 404 if project not found
- [ ] API returns 200 with project details

**Technical Notes:**
- Backend: GET /api/v1/projects/:id
- Frontend: Project detail page

---

### TM-004: Update Project
**As a** user  
**I want to** update a project  
**So that** I can modify project information  

**Acceptance Criteria:**
- [ ] User can update project name (unique validation)
- [ ] User can update project description
- [ ] User can update project status (active/archived)
- [ ] Updated timestamp is recorded
- [ ] Return 404 if project not found
- [ ] API returns 200 with updated project

**Technical Notes:**
- Backend: PATCH /api/v1/projects/:id
- Frontend: Edit project modal

---

### TM-005: Delete Project
**As a** user  
**I want to** delete a project  
**So that** I can remove projects I no longer need  

**Acceptance Criteria:**
- [ ] System confirms deletion before proceeding
- [ ] All tasks under the project are deleted
- [ ] Project is soft-deleted (marked as deleted)
- [ ] Return 404 if project not found
- [ ] API returns 204 on success

**Technical Notes:**
- Backend: DELETE /api/v1/projects/:id
- Frontend: Delete confirmation dialog

---

## Epic 2: Task Management

### TM-006: Create Task
**As a** user  
**I want to** create a task in a project  
**So that** I can track work items  

**Acceptance Criteria:**
- [ ] User can input task title (required, max 200 chars)
- [ ] User can input task description (optional, max 2000 chars)
- [ ] User can set priority (low, medium, high, urgent)
- [ ] User can set due date (optional, must be future date)
- [ ] User can assign labels (optional, multiple)
- [ ] System generates unique task ID
- [ ] Task is created with "todo" status
- [ ] Created timestamp is recorded
- [ ] Return 404 if project not found
- [ ] API returns 201 with task details

**Technical Notes:**
- Backend: POST /api/v1/projects/:projectId/tasks
- Frontend: Create task form in project view

---

### TM-007: List Tasks
**As a** user  
**I want to** view all tasks in a project  
**So that** I can see what needs to be done  

**Acceptance Criteria:**
- [ ] Display tasks with title, status, priority, due date
- [ ] Support sorting by: created date, due date, priority
- [ ] Support filtering by: status, priority, labels
- [ ] Support pagination (20 items per page)
- [ ] Overdue tasks are highlighted
- [ ] API returns 200 with tasks array

**Technical Notes:**
- Backend: GET /api/v1/projects/:projectId/tasks
- Frontend: Task list component with filters

---

### TM-008: Get Task Details
**As a** user  
**I want to** view a specific task  
**So that** I can see all its details  

**Acceptance Criteria:**
- [ ] Display task title, description, status, priority
- [ ] Display due date with overdue indicator
- [ ] Display labels assigned
- [ ] Display created and updated timestamps
- [ ] Display task dependencies (blocked by / blocking)
- [ ] Return 404 if task not found
- [ ] API returns 200 with task details

**Technical Notes:**
- Backend: GET /api/v1/tasks/:id
- Frontend: Task detail page/modal

---

### TM-009: Update Task
**As a** user  
**I want to** update a task  
**So that** I can modify task information  

**Acceptance Criteria:**
- [ ] User can update task title
- [ ] User can update task description
- [ ] User can update task status (todo, in-progress, done, cancelled)
- [ ] User can update task priority
- [ ] User can update due date
- [ ] User can update labels
- [ ] Updated timestamp is recorded
- [ ] Return 404 if task not found
- [ ] API returns 200 with updated task

**Technical Notes:**
- Backend: PATCH /api/v1/tasks/:id
- Frontend: Edit task modal

---

### TM-010: Delete Task
**As a** user  
**I want to** delete a task  
**So that** I can remove tasks I no longer need  

**Acceptance Criteria:**
- [ ] System confirms deletion before proceeding
- [ ] Task dependencies are removed
- [ ] Task is soft-deleted
- [ ] Return 404 if task not found
- [ ] API returns 204 on success

**Technical Notes:**
- Backend: DELETE /api/v1/tasks/:id
- Frontend: Delete confirmation dialog

---

## Epic 3: Task Dependencies

### TM-011: Add Task Dependency
**As a** user  
**I want to** mark a task as dependent on another task  
**So that** I can track task relationships  

**Acceptance Criteria:**
- [ ] User can specify a task that blocks current task
- [ ] System validates both tasks exist
- [ ] System prevents circular dependencies
- [ ] System prevents self-dependency
- [ ] System prevents duplicate dependencies
- [ ] Task status shows "blocked" if dependency is not done
- [ ] API returns 201 with dependency details

**Technical Notes:**
- Backend: POST /api/v1/tasks/:taskId/dependencies
- Frontend: Add dependency modal

---

### TM-012: Remove Task Dependency
**As a** user  
**I want to** remove a task dependency  
**So that** I can update task relationships  

**Acceptance Criteria:**
- [ ] User can remove a dependency from a task
- [ ] Dependency is removed from both tasks
- [ ] Return 404 if dependency not found
- [ ] API returns 204 on success

**Technical Notes:**
- Backend: DELETE /api/v1/tasks/:taskId/dependencies/:dependencyId
- Frontend: Dependency list with remove action

---

## Epic 4: Label Management

### TM-013: Create Label
**As a** user  
**I want to** create a label  
**So that** I can categorize tasks  

**Acceptance Criteria:**
- [ ] User can input label name (required, max 50 chars)
- [ ] User can select label color (hex color)
- [ ] System validates unique label name per user
- [ ] System generates unique label ID
- [ ] API returns 201 with label details

**Technical Notes:**
- Backend: POST /api/v1/labels
- Frontend: Label creation form

---

### TM-014: List Labels
**As a** user  
**I want to** view all my labels  
**So that** I can select labels for tasks  

**Acceptance Criteria:**
- [ ] Display all labels with name and color
- [ ] Show task count per label
- [ ] Support search by label name
- [ ] API returns 200 with labels array

**Technical Notes:**
- Backend: GET /api/v1/labels
- Frontend: Label picker component

---

### TM-015: Update Label
**As a** user  
**I want to** update a label  
**So that** I can modify label information  

**Acceptance Criteria:**
- [ ] User can update label name (unique validation)
- [ ] User can update label color
- [ ] Return 404 if label not found
- [ ] API returns 200 with updated label

**Technical Notes:**
- Backend: PATCH /api/v1/labels/:id
- Frontend: Edit label modal

---

### TM-016: Delete Label
**As a** user  
**I want to** delete a label  
**So that** I can remove labels I no longer need  

**Acceptance Criteria:**
- [ ] Label is removed from all tasks
- [ ] Label is soft-deleted
- [ ] Return 404 if label not found
- [ ] API returns 204 on success

**Technical Notes:**
- Backend: DELETE /api/v1/labels/:id
- Frontend: Delete confirmation dialog

---

## Epic 5: Search & Filtering

### TM-017: Global Task Search
**As a** user  
**I want to** search tasks across all projects  
**So that** I can find tasks quickly  

**Acceptance Criteria:**
- [ ] Search by task title (partial match, case-insensitive)
- [ ] Search by task description (partial match)
- [ ] Display results with project context
- [ ] Support pagination (20 items per page)
- [ ] Return empty array if no results
- [ ] API returns 200 with search results

**Technical Notes:**
- Backend: GET /api/v1/tasks/search?q=query
- Frontend: Global search component

---

### TM-018: Advanced Task Filtering
**As a** user  
**I want to** filter tasks with multiple criteria  
**So that** I can narrow down task list  

**Acceptance Criteria:**
- [ ] Filter by status (multiple selection)
- [ ] Filter by priority (multiple selection)
- [ ] Filter by labels (multiple selection)
- [ ] Filter by due date range
- [ ] Filter by overdue status
- [ ] Combine multiple filters (AND logic)
- [ ] API returns 200 with filtered tasks

**Technical Notes:**
- Backend: GET /api/v1/tasks?status=todo,in-progress&priority=high
- Frontend: Filter panel component

---

## Epic 6: Task Statistics

### TM-019: Project Statistics
**As a** user  
**I want to** see project statistics  
**So that** I can track project progress  

**Acceptance Criteria:**
- [ ] Display total task count
- [ ] Display tasks by status breakdown
- [ ] Display tasks by priority breakdown
- [ ] Display overdue task count
- [ ] Display completion percentage
- [ ] API returns 200 with statistics

**Technical Notes:**
- Backend: GET /api/v1/projects/:id/statistics
- Frontend: Project dashboard component

---

### TM-020: Dashboard Overview
**As a** user  
**I want to** see a dashboard overview  
**So that** I can quickly understand my work status  

**Acceptance Criteria:**
- [ ] Display total projects count
- [ ] Display total tasks count
- [ ] Display tasks due today
- [ ] Display tasks overdue
- [ ] Display recent tasks (last 5)
- [ ] API returns 200 with dashboard data

**Technical Notes:**
- Backend: GET /api/v1/dashboard
- Frontend: Dashboard page

---

## Technical Implementation Order

### Phase 1: Foundation (Cards TM-001 to TM-002)
1. TM-001: Create Project
2. TM-002: List Projects

### Phase 2: Core Task Management (Cards TM-003 to TM-010)
3. TM-003: Get Project Details
4. TM-006: Create Task
5. TM-007: List Tasks
6. TM-008: Get Task Details
7. TM-009: Update Task
8. TM-004: Update Project
9. TM-005: Delete Project
10. TM-010: Delete Task

### Phase 3: Advanced Features (Cards TM-011 to TM-020)
11. TM-011: Add Task Dependency
12. TM-012: Remove Task Dependency
13. TM-013: Create Label
14. TM-014: List Labels
15. TM-015: Update Label
16. TM-016: Delete Label
17. TM-017: Global Task Search
18. TM-018: Advanced Task Filtering
19. TM-019: Project Statistics
20. TM-020: Dashboard Overview

---

## API Endpoints Summary

### Projects
- POST /api/v1/projects
- GET /api/v1/projects
- GET /api/v1/projects/:id
- PATCH /api/v1/projects/:id
- DELETE /api/v1/projects/:id
- GET /api/v1/projects/:id/statistics

### Tasks
- POST /api/v1/projects/:projectId/tasks
- GET /api/v1/projects/:projectId/tasks
- GET /api/v1/tasks/:id
- PATCH /api/v1/tasks/:id
- DELETE /api/v1/tasks/:id
- GET /api/v1/tasks/search

### Dependencies
- POST /api/v1/tasks/:taskId/dependencies
- DELETE /api/v1/tasks/:taskId/dependencies/:dependencyId

### Labels
- POST /api/v1/labels
- GET /api/v1/labels
- PATCH /api/v1/labels/:id
- DELETE /api/v1/labels/:id

### Dashboard
- GET /api/v1/dashboard