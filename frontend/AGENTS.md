# Frontend AGENTS.md

## Project Overview

TaskFlow - Modern Task Management System Frontend - Next.js 14 with Bun runtime

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

design-system/       # Design system documentation
└── MASTER.md        # Design tokens & guidelines
```

## Tech Stack

- **Runtime**: Bun
- **Framework**: Next.js 14 (App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS
- **Font**: Plus Jakarta Sans (Google Fonts)
- **State**: Zustand + React Query
- **Testing**: Vitest + React Testing Library

## Design System

### Color Palette
- **Primary**: Teal (#0D9488)
- **CTA**: Orange (#F97316)
- **Background**: Teal 50 (#F0FDFA)
- **Text**: Gray 900 (#111827)

See `design-system/MASTER.md` for complete design tokens.

### UI Components

All components are in `components/ui/` and follow consistent patterns:

- **Button**: Primary, Secondary, CTA, Danger, Ghost variants
- **Card**: Default + hover states with shadows
- **Input**: With icons, labels, error states
- **Modal**: Accessible with animations
- **Badge**: Success, Warning, Info, Error, Neutral
- **Skeleton**: Loading states
- **EmptyState**: Empty data states

### Styling Conventions

```typescript
// Use Tailwind classes with design system tokens
<button className="btn-primary">        // Primary button
<div className="card-hover">            // Hoverable card
<input className="input">               // Styled input
<span className="badge-success">        // Success badge
```

### Custom Classes (defined in globals.css)

- `.btn`, `.btn-primary`, `.btn-secondary`, `.btn-cta`, `.btn-danger`, `.btn-ghost`
- `.card`, `.card-hover`
- `.input`, `.input-error`
- `.badge`, `.badge-success`, `.badge-warning`, `.badge-info`, `.badge-error`, `.badge-neutral`
- `.skeleton`

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
- Use design system colors (primary-600, cta-500, etc.)

### API Calls
- Use React Query for data fetching
- Define types in `lib/types/`

## Anti-Patterns (Avoid)

- ❌ useEffect for data fetching (use React Query)
- ❌ Prop drilling (use Context/Zustand)
- ❌ Inline styles (use Tailwind)
- ❌ Any types (use proper TypeScript)
- ❌ Console.log in production
- ❌ Emoji icons (use SVG instead)
- ❌ Missing cursor-pointer on interactive elements

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

## Accessibility Guidelines

- All interactive elements must have `cursor-pointer`
- Use semantic HTML elements
- Include ARIA labels where needed
- Support `prefers-reduced-motion`
- Minimum contrast ratio: 4.5:1
