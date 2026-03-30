# /component Command

Generate and update Solid UI components following project conventions.

## Usage

```bash
/component <description>
```

## Examples

### Create new components
```bash
/component create a tooltip that shows on hover
/component make a card for displaying server info
/component add a dropdown with search
/component generate a button with loading state
```

### Update existing components
```bash
/component update button to add loading state
/component fix combobox missing optionLabel prop
/component modify card to support dark mode
```

## How it works

1. **Parses natural language** - Understands "create", "make", "add", "update", "fix", etc.
2. **Detects component type** - Recognizes keywords like tooltip, button, card, dialog, etc.
3. **Chooses appropriate pattern:**
   - **Simple HTML Wrapper** - For basic styled elements (uses `cn()` utility)
   - **Kobalte Primitive Wrapper** - For interactive components (buttons, selects, etc.)
   - **CVA Variants** - For components with multiple visual states
4. **Generates files** in `apps/ui/src/components/ui/`
5. **Creates test stub** with basic vitest structure

## Generated structure

```
apps/ui/src/components/ui/
├── button.tsx         # Component file
└── button.test.tsx    # Test stub
```

## Conventions followed

- ✅ Kobalte primitives for interactive elements
- ✅ Tailwind CSS with theme variables
- ✅ TypeScript with proper generics
- ✅ Named exports following project patterns
- ✅ Import ordering from AGENTS.md
- ✅ No hardcoded colors

## Available Kobalte primitives

The tool can detect and use these primitives:
- Button, Select, Combobox, Checkbox, Dialog
- DropdownMenu, HoverCard, Label, Menubar
- Popover, RadioGroup, Slider, Switch, Tabs
- TextField, Tooltip

## Notes

- Components are created in `apps/ui/src/components/ui/`
- Test framework (vitest) not yet configured - tests generated but won't run
- Duplicate component names are detected and rejected
- Component names are auto-detected from keywords (tooltip → Tooltip)
