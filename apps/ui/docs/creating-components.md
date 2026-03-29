# Creating UI Components

This guide explains how to create new UI components using Kobalte primitives and Tailwind CSS.

## Overview

All UI components are built using one of three patterns:

1. **Simple HTML Wrapper** - For basic styled elements (Input, Label, Badge, Card)
2. **Kobalte Primitive Wrapper** - For complex interactive components (Button, Select, Combobox)
3. **CVA Variants** - For components with multiple visual states (Button, Badge)

## Pattern 1: Simple HTML Wrapper

Use this for basic elements that don't need complex interactions:

```typescript
import type { Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"
import { cn } from "~/lib/utils"

const Input: Component<ComponentProps<"input">> = (props) => {
  const [local, others] = splitProps(props, ["class", "type"])
  return (
    <input
      type={local.type || "text"}
      class={cn(
        "flex h-10 w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm ring-offset-background file:border-0 file:bg-transparent file:text-sm file:font-medium placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        local.class
      )}
      {...others}
    />
  )
}

export { Input }
export type { InputProps }
```

## Pattern 2: Kobalte Primitive Wrapper

Use this for complex interactive components:

```typescript
import type { JSX, ValidComponent } from "solid-js"
import { splitProps } from "solid-js"
import * as ButtonPrimitive from "@kobalte/core/button"
import type { PolymorphicProps } from "@kobalte/core/polymorphic"
import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"
import { cn } from "~/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground hover:bg-primary/90",
        outline: "border border-input hover:bg-accent hover:text-accent-foreground"
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 px-3"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default"
    }
  }
)

type ButtonProps<T extends ValidComponent = "button"> = ButtonPrimitive.ButtonRootProps<T> &
  VariantProps<typeof buttonVariants> & { class?: string | undefined; children?: JSX.Element }

const Button = <T extends ValidComponent = "button">(
  props: PolymorphicProps<T, ButtonProps<T>>
) => {
  const [local, others] = splitProps(props as ButtonProps, ["variant", "size", "class"])
  return (
    <ButtonPrimitive.Root
      class={cn(buttonVariants({ variant: local.variant, size: local.size }), local.class)}
      {...others}
    />
  )
}

export { Button, buttonVariants }
export type { ButtonProps }
```

## Pattern 3: CVA Variants

For components with multiple visual states:

```typescript
import type { Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"
import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"
import { cn } from "~/lib/utils"

const badgeVariants = cva(
  "inline-flex items-center rounded-md border px-2.5 py-0.5 text-xs font-semibold transition-colors",
  {
    variants: {
      variant: {
        default: "border-transparent bg-primary text-primary-foreground",
        secondary: "border-transparent bg-secondary text-secondary-foreground",
        outline: "text-foreground"
      }
    },
    defaultVariants: {
      variant: "default"
    }
  }
)

type BadgeProps = ComponentProps<"div"> & VariantProps<typeof badgeVariants>

const Badge: Component<BadgeProps> = (props) => {
  const [local, others] = splitProps(props, ["class", "variant"])
  return (
    <div
      class={cn(badgeVariants({ variant: local.variant }), local.class)}
      {...others}
    />
  )
}

export { Badge, badgeVariants }
export type { BadgeProps }
```

## Step-by-Step Process

1. **Determine the pattern** - Does it need Kobalte? Does it have variants?
2. **Create the file** - `apps/ui/src/components/ui/<component-name>.tsx`
3. **Copy the template** - Use the appropriate pattern above
4. **Import required primitives** - See available Kobalte primitives below
5. **Define styling** - Use Tailwind classes with theme variables
6. **Export properly** - Export component, variants (if CVA), and types
7. **Test it** - Import in a route and verify it works

## Available Kobalte Primitives

- `@kobalte/core/button` - Buttons, actions
- `@kobalte/core/select` - Dropdown selects
- `@kobalte/core/combobox` - Autocomplete/search selects
- `@kobalte/core/checkbox` - Checkboxes
- `@kobalte/core/dialog` - Modal dialogs
- `@kobalte/core/dropdown-menu` - Context menus
- `@kobalte/core/hover-card` - Hover popovers
- `@kobalte/core/label` - Form labels (Note: we use plain HTML `<label>`)
- `@kobalte/core/menubar` - Menu bars
- `@kobalte/core/popover` - Floating popovers
- `@kobalte/core/radio-group` - Radio buttons
- `@kobalte/core/slider` - Range sliders
- `@kobalte/core/switch` - Toggle switches
- `@kobalte/core/tabs` - Tab panels
- `@kobalte/core/text-field` - Form inputs (Note: we use plain HTML `<input>`)
- `@kobalte/core/tooltip` - Tooltips

## Styling Guidelines

### DO:
- Use theme CSS variables: `bg-background`, `text-foreground`, `text-muted-foreground`
- Use the `cn()` utility for class merging
- Trust Kobalte primitive styles - only override when necessary
- Follow existing component patterns

### DON'T:
- Use hardcoded colors like `bg-[#0d0d0d]` or `text-white/90`
- Modify CSS theme in `src/app.css`
- Use `npx solidui-cli` or `npx shadcn` (these are broken/React tools)
- Override Kobalte styles unnecessarily

## Example: Creating a New Component

Let's create a `Textarea` component:

```typescript
// apps/ui/src/components/ui/textarea.tsx
import type { Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"
import { cn } from "~/lib/utils"

const Textarea: Component<ComponentProps<"textarea">> = (props) => {
  const [local, others] = splitProps(props, ["class"])
  return (
    <textarea
      class={cn(
        "flex min-h-[80px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
        local.class
      )}
      {...others}
    />
  )
}

export { Textarea }
```

## Testing Your Component

1. Import it in a route file:
   ```typescript
   import { Textarea } from "~/components/ui/textarea"
   ```

2. Use it in your component:
   ```tsx
   <Textarea placeholder="Enter your message" />
   ```

3. Run type checking:
   ```bash
   cd apps/ui && pnpm check
   ```

4. Build to ensure no bundling issues:
   ```bash
   cd apps/ui && pnpm build
   ```

## Reference Components

Look at these existing components for examples:

- **Simple HTML**: `input.tsx`, `label.tsx`, `card.tsx`
- **Kobalte Primitive**: `button.tsx`, `select.tsx`, `combobox.tsx`
- **CVA Variants**: `button.tsx`, `badge.tsx`
