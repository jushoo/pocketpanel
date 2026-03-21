# Agent Guidelines for PocketPanel Web

## Component Usage

**DO NOT** create components manually that can be installed via shadcn-svelte CLI. Always use the CLI to install components:

```bash
npx shadcn-svelte@latest add <component-name>
```

Available components can be found at https://shadcn-svelte.com

## Styling Rules

1. **DO NOT modify the CSS theme** in `src/routes/layout.css`
   - The color theme is managed by shadcn-svelte
   - Use CSS variables like `--background`, `--foreground`, `--muted`, etc.
   - Theme switching will be added later - keep light/dark mode support intact

2. **ALWAYS respect the existing color theme**
   - Use Tailwind classes like `bg-background`, `text-foreground`, `text-muted-foreground`
   - Do not use hardcoded colors like `bg-[#0d0d0d]` or `text-white/90`
   - Let shadcn-svelte components use their default styling

## UI Design Principles

Design UIs that are:

- **Clean** - Minimal visual noise, generous whitespace
- **Minimal** - Only essential elements, no decorative fluff
- **Basic** - Simple layouts, straightforward interactions
- **Ollama-style** - Centered content, subtle borders, muted text hierarchy

### Examples of Good Design

```svelte
<!-- Good: Uses theme variables, clean layout -->
<div class="flex min-h-screen items-center justify-center bg-background">
	<div class="w-full max-w-sm px-6">
		<h1 class="text-xl font-medium text-foreground">Title</h1>
		<p class="text-sm text-muted-foreground">Subtitle</p>
		<Input class="h-11" />
		<Button class="w-full">Action</Button>
	</div>
</div>
```

```svelte
<!-- Bad: Hardcoded colors, custom styling overrides -->
<div class="min-h-screen bg-[#0d0d0d]">
	<h1 class="text-white/90">Title</h1>
	<Input class="border-white/10 bg-white/[0.03] text-white/90" />
	<Button class="bg-white/90 text-black">Action</Button>
</div>
```

## Summary

- Use shadcn-svelte CLI for components
- Never modify layout.css theme colors
- Use theme CSS variables for colors
- Keep designs clean, minimal, and ollama-inspired
- Trust the shadcn-svelte default component styling
