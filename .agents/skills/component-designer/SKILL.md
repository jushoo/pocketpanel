---
name: component-designer
description: Generates and updates Solid UI components using Kobalte primitives and Tailwind CSS. Creates components in apps/ui/src/components/ui/ following project conventions.
category: ui
tools:
  - name: generate-component
    description: Generate a new Solid UI component from natural language description
    parameters:
      - name: description
        type: string
        description: Natural language description of the component to create (e.g., "create a dropdown with search")
      - name: component_name
        type: string
        description: The component name in PascalCase (optional - will be inferred from description)
  - name: update-component
    description: Update an existing component with specified changes
    parameters:
      - name: component_path
        type: string
        description: Path to the component file relative to apps/ui/src/components/ui/
      - name: changes
        type: string
        description: Natural language description of changes to make
---

# Component Designer

Generate and update Solid UI components following project conventions from AGENTS.md.

## Patterns

### Pattern 1: Simple HTML Wrapper
For basic styled elements without interaction (Card, Input, Label, Badge)
- Use `Component<ComponentProps<"element">>`
- Use `splitProps` to extract `class` prop
- Apply Tailwind classes via `cn()` utility
- Spread remaining props

### Pattern 2: Kobalte Primitive Wrapper
For complex interactive components (Button, Select, Combobox, Dialog)
- Import from `@kobalte/core/<primitive>`
- Use `PolymorphicProps` for type safety
- Wrap primitive with custom styling via `cn()`
- Support variants using `class-variance-authority`

### Pattern 3: CVA Variants
For components with multiple visual states
- Define variants with `cva()`
- Support size, color, and state variants
- Type-safe props using `VariantProps`

## Conventions

**Import Order:**
```typescript
// 1. Solid imports
import { splitProps } from "solid-js"
import type { Component, ComponentProps, ValidComponent, JSX } from "solid-js"

// 2. Kobalte imports
import * as Primitive from "@kobalte/core/<primitive>"
import type { PolymorphicProps } from "@kobalte/core/polymorphic"

// 3. Utility imports
import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"
import { cn } from "~/lib/utils"
```

**Naming:**
- Components: PascalCase (`Button.tsx`)
- Exports: Named exports for main component
- Compound components: Export all parts (e.g., `Card`, `CardHeader`, `CardContent`)

**Styling Rules:**
- NEVER modify CSS theme variables in `src/app.css`
- Use theme CSS variables: `bg-background`, `text-foreground`, `text-muted-foreground`
- NO hardcoded colors - use Tailwind classes
- Trust Kobalte primitive defaults - don't override unnecessarily

**Design Principles:**
- Clean: Minimal visual noise, generous whitespace
- Minimal: Only essential elements
- Basic: Simple layouts, straightforward interactions
- Ollama-style: Centered content, subtle borders, muted text hierarchy

## Available Kobalte Primitives

- `@kobalte/core/button` - Buttons, actions
- `@kobalte/core/select` - Dropdown selects
- `@kobalte/core/combobox` - Autocomplete/search selects
- `@kobalte/core/checkbox` - Checkboxes
- `@kobalte/core/dialog` - Modal dialogs
- `@kobalte/core/dropdown-menu` - Context menus
- `@kobalte/core/hover-card` - Hover popovers
- `@kobalte/core/label` - Form labels
- `@kobalte/core/menubar` - Menu bars
- `@kobalte/core/popover` - Floating popovers
- `@kobalte/core/radio-group` - Radio buttons
- `@kobalte/core/slider` - Range sliders
- `@kobalte/core/switch` - Toggle switches
- `@kobalte/core/tabs` - Tab panels
- `@kobalte/core/text-field` - Form inputs
- `@kobalte/core/tooltip` - Tooltips

## Component Generation Workflow

### 1. Analyze the Request
- Determine component type: simple, interactive, or variant-based
- Identify required Kobalte primitives (if any)
- Check for compound components (Header/Content/Footer pattern)

### 2. Read Existing Patterns
Before creating a new component:
- Check `_TEMPLATE.tsx` for base patterns
- Look at similar existing components (`card.tsx`, `button.tsx`, etc.)
- Follow established naming and structure conventions

### 3. Generate Component
Create file at `apps/ui/src/components/ui/<name>.tsx` with:
- Proper imports following import order
- Type-safe props using Solid and Kobalte types
- Tailwind classes using theme variables
- Compound components if applicable
- Named exports for all public parts

### 4. Generate Test Stub
Create file at `apps/ui/src/components/ui/<name>.test.tsx` with:
- Basic vitest test structure
- Component rendering test
- Accessibility checks if interactive

**Note:** The project doesn't have vitest configured yet. The test file will be generated but won't run until vitest is set up.

### 5. Update Component (if requested)
When updating an existing component:
- Read the current component file
- Understand its structure and exports
- Apply changes while preserving existing API
- Maintain backward compatibility when possible
- Update tests to cover new functionality

## Examples

### Generate Request
**Input:** "create a card component for displaying server info"

**Action:**
1. Check `card.tsx` - already exists
2. Suggest creating a specialized variant or updating existing

### Generate Request
**Input:** "create a tooltip that shows on hover with dark background"

**Action:**
1. Use Pattern 2 (Kobalte wrapper)
2. Import `@kobalte/core/tooltip`
3. Wrap with custom dark styling
4. Export `Tooltip`, `TooltipTrigger`, `TooltipContent`

### Update Request
**Input:** "update button to add a loading state with spinner"

**Action:**
1. Read existing `button.tsx`
2. Add `isLoading` prop
3. Conditionally render spinner icon
4. Update `buttonVariants` to include loading styles
5. Update tests for loading state

## Output Format

When generating a component, provide:
1. **Component code** with proper TypeScript types
2. **Test stub** with basic tests
3. **Usage example** showing how to import and use
4. **Pattern notes** explaining which pattern was used and why

## Validation Checklist

Before finalizing a component:
- [ ] Follows import ordering from AGENTS.md
- [ ] Uses theme CSS variables (no hardcoded colors)
- [ ] Proper TypeScript types with generics if polymorphic
- [ ] Kobalte primitives properly wrapped if interactive
- [ ] Tailwind classes use `cn()` utility
- [ ] Named exports for all public parts
- [ ] Test file generated with basic structure
- [ ] No emojis unless explicitly requested
- [ ] Follows Solid/SolidStart conventions from AGENTS.md
