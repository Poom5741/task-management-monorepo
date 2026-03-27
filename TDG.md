# TDG Configuration

## Project Information
- Language: Go (backend), TypeScript (frontend)
- Framework: Gin (backend), Next.js 14 (frontend)
- Test Framework: testify (backend), Vitest + React Testing Library (frontend)

## Build Command
- Backend: `cd backend && make build`
- Frontend: `bun run --cwd frontend build`

## Test Command
- Backend: `cd backend && go test ./...`
- Frontend: `bun run --cwd frontend test`

## Single Test Command
- Backend: `cd backend && go test -v -run TestName ./path/to/package`
- Frontend: `bun run --cwd frontend test -- TestName`

## Coverage Command
- Backend: `cd backend && go test -cover ./...`
- Frontend: `bun run --cwd frontend test:coverage`

## Test File Patterns
- Backend: `*_test.go` in `backend/internal/`
- Frontend: `*.test.ts`, `*.test.tsx` in `frontend/`