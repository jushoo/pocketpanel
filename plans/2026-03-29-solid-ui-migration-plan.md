# Migration: shadcn-solid to solid-ui

**Created**: 2026-03-29  
**Status**: Planning

## Overview

Migrate the `apps/ui` project from shadcn-solid component library to solid-ui (https://www.solid-ui.com/). This migration will align the project with the more actively maintained solid-ui ecosystem while maintaining the current Tailwind CSS v4 styling system.

## Context

### Current State
- **Project**: `apps/ui` - A SolidStart application
- **UI Library**: shadcn-solid (unofficial port of shadcn/ui)
- **Schema**: `https://shadcn-solid.com/schema.json`
- **Components**: Button, Card (suite), Input, Label, Badge
- **Styling**: Tailwind CSS v4 with tw-animate-css
- **Utilities**: class-variance-authority, tailwind-merge, clsx

### Desired State
- **UI Library**: solid-ui (more active, better maintained)
- **Schema**: `https://solid-ui.com/schema.json`
- **Same components** but from solid-ui registry
- **Preserved styling** - maintain current Tailwind v4 setup
- **Updated imports** - components may have slightly different APIs

### Motivation
1. **Active maintenance** - solid-ui is more actively maintained
2. **Better documentation** - comprehensive docs at solid-ui.com
3. **More components** - access to additional components if needed
4. **Community alignment** - following the standard Solid ecosystem

## Scope

### In Scope
- [ ] Update `components.json` configuration for solid-ui
- [ ] Install solid-ui CLI and dependencies
- [ ] Migrate existing components (Button, Card, Input, Label, Badge)
- [ ] Update all route files to use new component imports
- [ ] Verify TypeScript types are correct
- [ ] Ensure all tests pass
- [ ] Verify build succeeds

### Out of Scope
- Adding new components (future enhancement)
- Changing styling/theme system (preserve current Tailwind v4 setup)
- Modifying business logic in routes
- Backend API changes

## Current Component Inventory

| Component | Current Location | Import Path | Solid-UI Equivalent |
|-----------|-------------------|-------------|---------------------|
| Button | `~/components/ui/button/button.tsx` | `~/components/ui/button` | `button` |
| Card | `~/components/ui/card/card.tsx` | `~/components/ui/card` | `card` |
| CardHeader | `~/components/ui/card/card-header.tsx` | `~/components/ui/card` | Part of card |
| CardTitle | `~/components/ui/card/card-title.tsx` | `~/components/ui/card` | Part of card |
| CardDescription | `~/components/ui/card/card-description.tsx` | `~/components/ui/card` | Part of card |
| CardContent | `~/components/ui/card/card-content.tsx` | `~/components/ui/card` | Part of card |
| CardFooter | `~/components/ui/card/card-footer.tsx` | `~/components/ui/card` | Part of card |
| Input | `~/components/ui/input/input.tsx` | `~/components/ui/input` | `input` |
| Label | `~/components/ui/label/label.tsx` | `~/components/ui/label` | `label` |
| Badge | `~/components/ui/badge/badge.tsx` | `~/components/ui/badge` | `badge` |

## Files Using UI Components

1. `apps/ui/src/routes/index.tsx` - Login page (Button, Input, Label)
2. `apps/ui/src/routes/dashboard/index.tsx` - Dashboard (Button, Card suite, Badge)
3. `apps/ui/src/routes/dashboard/servers/new.tsx` - Create server (Button, Input, Label, Card suite)
4. `apps/ui/src/routes/dashboard/servers/[id].tsx` - Server detail (Button, Card suite, Badge, Input, Label)
5. `apps/ui/src/components/ModeToggle.tsx` - Theme toggle (Button)

## Implementation Plan

### Step 1: Preparation & Configuration
**What**: Update project configuration for solid-ui
**Files**:
- `apps/ui/components.json`
- `apps/ui/package.json`
**Tasks**:
- [ ] Update `components.json` schema URL to `https://solid-ui.com/schema.json`
- [ ] Install solid-ui CLI: `npx solid-ui@latest init`
- [ ] Verify `tailwind.config` compatibility (should work with v4)
- [ ] Keep existing `app.css` styling - no changes needed
**Acceptance Criteria**:
- `npx solid-ui@latest add --help` works
- CLI can add components without errors

### Step 2: Install Solid-UI Components
**What**: Add solid-ui components one by one to replace existing ones
**Files**:
- New: `apps/ui/src/components/ui/button.tsx`
- New: `apps/ui/src/components/ui/card.tsx`
- New: `apps/ui/src/components/ui/input.tsx`
- New: `apps/ui/src/components/ui/label.tsx`
- New: `apps/ui/src/components/ui/badge.tsx`
- Delete: Old component directories
**Tasks**:
- [ ] Run `npx solid-ui@latest add button`
- [ ] Run `npx solid-ui@latest add card`
- [ ] Run `npx solid-ui@latest add input`
- [ ] Run `npx solid-ui@latest add label`
- [ ] Run `npx solid-ui@latest add badge`
- [ ] Review each generated component for compatibility
- [ ] Delete old component folders: `button/`, `card/`, `input/`, `label/`, `badge/`
**Acceptance Criteria**:
- All components installed in `~/components/ui/` as single files
- Components compile without TypeScript errors
- Visual appearance matches current styling

### Step 3: Update Component Imports
**What**: Update all files that import from old component paths
**Files**:
- `apps/ui/src/routes/index.tsx`
- `apps/ui/src/routes/dashboard/index.tsx`
- `apps/ui/src/routes/dashboard/servers/new.tsx`
- `apps/ui/src/routes/dashboard/servers/[id].tsx`
- `apps/ui/src/components/ModeToggle.tsx`
**Tasks**:
- [ ] Update imports from `~/components/ui/button` to `~/components/ui/button`
- [ ] Update imports from `~/components/ui/card` to `~/components/ui/card`
- [ ] Update imports from `~/components/ui/input` to `~/components/ui/input`
- [ ] Update imports from `~/components/ui/label` to `~/components/ui/label`
- [ ] Update imports from `~/components/ui/badge` to `~/components/ui/badge`
- [ ] Check for any API changes (solid-ui may have slightly different props)
**Acceptance Criteria**:
- All imports resolve correctly
- No "module not found" errors
- Components render correctly in UI

### Step 4: Verify & Test
**What**: Ensure everything works correctly
**Files**:
- All route files
- Component files
**Tasks**:
- [ ] Run `pnpm check` - TypeScript check
- [ ] Run `pnpm lint` - Linting
- [ ] Run `pnpm dev` and verify UI visually
- [ ] Test login page functionality
- [ ] Test dashboard server creation
- [ ] Test dark mode toggle
- [ ] Verify all component variants work
**Acceptance Criteria**:
- No TypeScript errors
- No linting errors
- UI looks identical to before migration
- All user flows work correctly

### Step 5: Cleanup
**What**: Remove any remaining shadcn-solid artifacts
**Files**:
- Check for any residual configuration
**Tasks**:
- [ ] Verify no old shadcn-solid references in codebase
- [ ] Check `node_modules` is clean (fresh install if needed)
- [ ] Update any documentation/comments referencing shadcn-solid
**Acceptance Criteria**:
- No references to shadcn-solid in codebase
- Clean dependency tree

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Solid-ui components have different APIs | Medium | Review each component's API before migration; adjust props as needed |
| Tailwind v4 compatibility issues | Medium | Test early; solid-ui should support v4 but verify |
| Build size increase | Low | Monitor bundle size; tree-shaking should minimize impact |
| Breaking changes in route logic | Low | Only changing imports, not logic; test all routes |
| Loss of custom styling | Low | Solid-ui should respect CSS variables; verify visual match |

## Dependencies

- solid-ui CLI working in the project
- Tailwind CSS v4 compatibility confirmed
- Access to install new npm packages

## Testing Strategy

1. **Static Analysis**
   - TypeScript check: `pnpm check`
   - Linting: `pnpm lint`
   - Build: `pnpm build`

2. **Visual Regression**
   - Compare login page before/after
   - Compare dashboard before/after
   - Verify dark mode works
   - Check responsive behavior

3. **Functional Testing**
   - Login flow works
   - Server creation works
   - Navigation works
   - All buttons clickable
   - All inputs functional

## Rollback Plan

If issues arise:
1. Restore old `components.json` from git
2. Reinstall old components: `git checkout apps/ui/src/components/ui/`
3. Revert import changes in route files
4. Remove solid-ui from dependencies

## Checklist

- [ ] Step 1: Configuration complete
- [ ] Step 2: All components installed
- [ ] Step 3: All imports updated
- [ ] Step 4: All tests passing
- [ ] Step 5: Cleanup complete
- [ ] Visual verification complete
- [ ] Functional testing complete
- [ ] Ready for review

## Notes

- Keep the current Tailwind CSS v4 setup with tw-animate-css
- The `app.css` file should remain unchanged as it defines the theme
- `cn()` utility in `~/lib/utils.ts` should still work with solid-ui
- Component structure changes from multi-file (button/button.tsx, button/index.ts) to single file (button.tsx)
