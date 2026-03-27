# Frontend TDG.md

## Test Configuration

### Build Command
```bash
bun run build
```

### Test Command (Full Suite)
```bash
bun run test
```

### Single Test Command
```bash
bun run test ComponentName.test.tsx
```

### Coverage Command
```bash
bun run test:coverage
```

### Test File Patterns

#### Component Tests
- Location: `components/{category}/{ComponentName}.test.tsx`
- Pattern: Test file next to component

#### Page Tests
- Location: `app/{route}/page.test.tsx`
- Pattern: Test file next to page

#### Utility Tests
- Location: `lib/{util}.test.ts`
- Pattern: Test file next to utility

### Coverage Threshold
- Minimum: 70%
- Target: 85%

### Test Setup
Tests use:
- Vitest as test runner
- jsdom for DOM simulation
- @testing-library/react for component testing

### Mock Patterns

#### API Mock
```typescript
vi.mock('@/lib/api', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}))
```

#### Next.js Router Mock
```typescript
vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
    back: vi.fn(),
  }),
}))
```

### Pre-commit Hooks
Run before every commit:
1. `bun run lint` - Run linter
2. `bun run test` - Run tests
3. `bun run build` - Build check