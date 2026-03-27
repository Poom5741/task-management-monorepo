# Frontend AGENTS.md

## Project Overview

Task Management System Frontend - Next.js 14 with Bun runtime

## Architecture

This project follows **Feature-Based Architecture**:

```
app/                 # Next.js App Router pages
├── (dashboard)/     # Route group for dashboard
├── projects/        # Projects feature pages
├── tasks/           # Tasks feature pages
└── labels/          # Labels feature pages

components/
├── ui/              # Reusable UI components
└── layout/          # Layout components

lib/
├── api.ts           # API client
├── types/           # TypeScript types
└── utils.ts         # Utility functions
```

## Tech Stack

- **Runtime**: Bun
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **State**: Zustand + React Query
- **Testing**: Vitest + React Testing Library

## Testing Patterns

### Component Tests
```typescript
import { render, screen, fireEvent } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

describe('ProjectCard', () => {
  it('should render project details', () => {
    render(
      <QueryClientProvider client={new QueryClient()}>
        <ProjectCard project={mockProject} />
      </QueryClientProvider>
    )
    expect(screen.getByText('Project Name')).toBeInTheDocument()
  })
})
```

### Test Categories
1. **Rendering** - Component renders correctly
2. **User Interaction** - Click, type, submit
3. **State Changes** - Loading, error, success
4. **Accessibility** - ARIA attributes

## Code Conventions

### Component Structure
```typescript
// 1. Imports
import { useState } from 'react'

// 2. Types
interface Props {
  title: string
}

// 3. Component
export function Component({ title }: Props) {
  // hooks
  // handlers
  // render
  return <div>{title}</div>
}
```

### Styling
- Use Tailwind CSS classes
- Use `cn()` utility for conditional classes
- Mobile-first responsive design

### API Calls
- Use React Query for data fetching
- Define types in `lib/types/`

## Anti-Patterns (Avoid)

- ❌ useEffect for data fetching (use React Query)
- ❌ Prop drilling (use Context/Zustand)
- ❌ Inline styles (use Tailwind)
- ❌ Any types (use proper TypeScript)
- ❌ Console.log in production

## Commands

```bash
# Development
bun run dev

# Build
bun run build

# Test
bun run test

# Test with UI
bun run test:ui

# Coverage
bun run test:coverage

# Lint
bun run lint
```

## File Naming

- Components: PascalCase (e.g., `ProjectCard.tsx`)
- Utilities: camelCase (e.g., `formatDate.ts`)
- Pages: lowercase with dashes (e.g., `project-list.tsx`)