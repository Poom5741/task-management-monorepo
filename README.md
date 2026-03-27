# Task Management System - Monorepo

A full-stack Task Management application with Go backend (hexagonal architecture) and Next.js frontend (Bun).

## Project Structure

```
.
├── backend/                 # Go backend (hexagonal architecture)
│   ├── cmd/
│   │   └── server/         # Application entry point
│   ├── internal/
│   │   ├── domain/         # Domain entities & business rules
│   │   │   ├── project/
│   │   │   ├── task/
│   │   │   ├── label/
│   │   │   └── dependency/
│   │   ├── handler/        # HTTP handlers (delivery layer)
│   │   ├── usecase/        # Business logic (application layer)
│   │   └── storage/        # Data access (infrastructure layer)
│   │       ├── postgres/
│   │       └── redis/
│   ├── pkg/                # Shared packages
│   │   ├── config/
│   │   ├── logger/
│   │   └── validator/
│   ├── Makefile
│   └── go.mod
│
├── frontend/               # Next.js frontend (Bun)
│   ├── app/                # App router pages
│   ├── components/          # React components
│   ├── lib/                 # Utilities
│   ├── public/              # Static assets
│   └── package.json
│
├── docs/                   # Documentation
│   └── card.md            # Feature cards (Jira-style)
│
├── package.json            # Root package.json (workspaces)
└── README.md
```

## Tech Stack

### Backend
- **Language**: Go 1.22+
- **Architecture**: Hexagonal (Ports & Adapters)
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Testing**: Testcontainers-Go

### Frontend
- **Runtime**: Bun
- **Framework**: Next.js 14+
- **Testing**: Vitest + React Testing Library

## Development

### Prerequisites
- Go 1.22+
- Bun
- PostgreSQL
- Redis

### Setup

```bash
# Install all dependencies
bun run install:all

# Copy environment variables
cp backend/.env.example backend/.env

# Run development servers
bun run dev
```

### Backend Commands

```bash
cd backend

# Run server
make run

# Run tests
make test

# Run tests with coverage
make test-coverage

# Build binary
make build
```

### Frontend Commands

```bash
bun run frontend dev      # Development server
bun run frontend build    # Production build
bun run frontend test     # Run tests
bun run frontend lint     # Run linter
```

## TDG Workflow

This project follows Test-Driven Generation (TDG) workflow with Red-Green-Refactor cycles.

### Commit Message Format

```
red: test spec for <feature> (#issue-number)
green: implement <feature> (#issue-number)
refactor: <description> (#issue-number)
```

### Development Phases

1. **RED Phase**: Write failing tests first
2. **GREEN Phase**: Implement code to pass tests
3. **REFACTOR Phase**: Optimize and clean up code

## API Endpoints

See `docs/card.md` for complete API documentation.

## License

MIT