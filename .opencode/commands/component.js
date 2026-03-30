#!/usr/bin/env node

/**
 * /component command - Generate and update Solid UI components
 * 
 * Usage:
 *   /component create a card for displaying server info
 *   /component update button to add loading state
 *   /component fix combobox missing optionLabel prop
 * 
 * Components are created in: apps/ui/src/components/ui/
 */

const fs = require('fs');
const path = require('path');

const UI_DIR = path.join(process.cwd(), 'apps/ui/src/components/ui');
const SKILL_PATH = path.join(process.cwd(), '.agents/skills/component-designer/SKILL.md');

// Pattern templates based on _TEMPLATE.tsx
const PATTERNS = {
  simple: (name, description) => `import type { Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"

import { cn } from "~/lib/utils"

/**
 * ${name} component
 * ${description}
 */
const ${name}: Component<ComponentProps<"div">> = (props) => {
  const [local, others] = splitProps(props, ["class"])
  return (
    <div
      class={cn(
        "rounded-lg border bg-card text-card-foreground",
        local.class
      )}
      {...others}
    />
  )
}

export { ${name} }`,

  kobalte: (name, primitive, description) => `import type { ValidComponent, JSX } from "solid-js"
import { splitProps } from "solid-js"

import * as ${primitive}Primitive from "@kobalte/core/${primitive.toLowerCase()}"
import type { PolymorphicProps } from "@kobalte/core/polymorphic"

import { cn } from "~/lib/utils"

type ${name}Props<T extends ValidComponent = "button"> = 
  ${primitive}Primitive.${primitive}RootProps<T> & { 
    class?: string
    children?: JSX.Element 
  }

/**
 * ${name} component
 * ${description}
 */
const ${name} = <T extends ValidComponent = "button">(
  props: PolymorphicProps<T, ${name}Props<T>>
) => {
  const [local, others] = splitProps(props as ${name}Props, ["class", "children"])
  return (
    <${primitive}Primitive.Root
      class={cn(
        "inline-flex items-center justify-center",
        local.class
      )}
      {...others}
    >
      {local.children}
    </${primitive}Primitive.Root>
  )
}

export { ${name} }
export type { ${name}Props }`,

  cva: (name, description) => `import type { Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"
import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"

import { cn } from "~/lib/utils"

const ${name.toLowerCase()}Variants = cva(
  "inline-flex items-center justify-center rounded-md",
  {
    variants: {
      variant: {
        default: "bg-primary text-primary-foreground",
        outline: "border border-input bg-background",
        ghost: "hover:bg-accent hover:text-accent-foreground"
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 px-3",
        lg: "h-11 px-8"
      }
    },
    defaultVariants: {
      variant: "default",
      size: "default"
    }
  }
)

type ${name}Props = ComponentProps<"div"> &
  VariantProps<typeof ${name.toLowerCase()}Variants>

/**
 * ${name} component
 * ${description}
 */
const ${name}: Component<${name}Props> = (props) => {
  const [local, others] = splitProps(props, ["class", "variant", "size"])
  return (
    <div
      class={cn(
        ${name.toLowerCase()}Variants({ variant: local.variant, size: local.size }),
        local.class
      )}
      {...others}
    />
  )
}

export { ${name}, ${name.toLowerCase()}Variants }
export type { ${name}Props }`
};

// Test template
const TEST_TEMPLATE = (name) => `import { describe, it, expect } from "vitest"
import { render } from "@solidjs/testing-library"
import { ${name} } from "./${name.toLowerCase()}"

describe("${name}", () => {
  it("should render without crashing", () => {
    const { container } = render(() => <${name} />)
    expect(container).toBeTruthy()
  })

  it("should apply custom class", () => {
    const { container } = render(() => <${name} class="custom-class" />)
    expect(container.firstChild).toHaveClass("custom-class")
  })
})`;

// Parse natural language request
function parseRequest(input) {
  const normalized = input.toLowerCase().trim();
  
  // Check for update requests
  const updateMatch = normalized.match(/^(update|fix|modify|change)\s+(?:the\s+)?(\w+)(?:\s+component)?\s+(?:to\s+)?(.+)$/i);
  if (updateMatch) {
    return {
      type: 'update',
      componentName: updateMatch[2],
      description: updateMatch[3],
      raw: input
    };
  }
  
  // Check for create requests
  const createMatch = normalized.match(/^(create|make|generate|add)\s+(?:a|an)?\s*(.+)$/i);
  if (createMatch) {
    // Try to extract component name from description
    const description = createMatch[2];
    const nameMatch = description.match(/(\w+)\s+(?:component|card|button|input|modal|dialog|tooltip|dropdown|select)/i);
    const componentName = nameMatch ? toPascalCase(nameMatch[1]) : null;
    
    return {
      type: 'create',
      componentName,
      description,
      raw: input
    };
  }
  
  // Default to create if unclear
  return {
    type: 'create',
    componentName: null,
    description: input,
    raw: input
  };
}

// Convert to PascalCase
function toPascalCase(str) {
  return str
    .split(/[-_\s]+/)
    .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
    .join('');
}

// Detect appropriate pattern
function detectPattern(description) {
  const normalized = description.toLowerCase();
  
  // Check for interactive keywords
  const interactiveKeywords = ['button', 'click', 'hover', 'focus', 'select', 'dropdown', 'modal', 'dialog', 'toggle', 'checkbox', 'radio', 'slider', 'tooltip', 'popover', 'combobox'];
  const isInteractive = interactiveKeywords.some(kw => normalized.includes(kw));
  
  // Check for variant keywords
  const variantKeywords = ['variant', 'size', 'color', 'style', 'type', 'primary', 'secondary', 'outline', 'ghost'];
  const hasVariants = variantKeywords.some(kw => normalized.includes(kw));
  
  if (hasVariants) return 'cva';
  if (isInteractive) return 'kobalte';
  return 'simple';
}

// Extract component type from keywords
const COMPONENT_TYPES = {
  tooltip: 'Tooltip',
  button: 'Button',
  card: 'Card',
  input: 'Input',
  select: 'Select',
  dropdown: 'Dropdown',
  combobox: 'Combobox',
  dialog: 'Dialog',
  modal: 'Modal',
  checkbox: 'Checkbox',
  radio: 'RadioGroup',
  slider: 'Slider',
  switch: 'Switch',
  tabs: 'Tabs',
  badge: 'Badge',
  label: 'Label',
  popover: 'Popover',
  menu: 'Menu',
  menubar: 'Menubar',
  hover: 'HoverCard',
  textfield: 'TextField'
};

// Generate component name from description
function generateName(description) {
  const normalized = description.toLowerCase();
  
  // First, check for explicit component type keywords
  for (const [keyword, type] of Object.entries(COMPONENT_TYPES)) {
    if (normalized.includes(keyword)) {
      return type;
    }
  }
  
  // Otherwise, extract meaningful words
  const words = description
    .toLowerCase()
    .replace(/(?:create|make|generate|add|a|an|the|for|with|component|that|shows|displays|on|hover|click)/g, ' ')
    .trim()
    .split(/[\s-_]+/)
    .filter(w => w.length > 2 && !['and', 'the', 'for', 'with'].includes(w));
  
  if (words.length === 0) return 'Component';
  
  // Take first meaningful word or combine
  if (words.length === 1) {
    return toPascalCase(words[0]);
  }
  
  // For multi-word names, use first 2 words
  return toPascalCase(words.slice(0, 2).join('-'));
}

// Read existing component
function readComponent(componentName) {
  const fileName = componentName.toLowerCase() + '.tsx';
  const filePath = path.join(UI_DIR, fileName);
  
  if (!fs.existsSync(filePath)) {
    return null;
  }
  
  return {
    path: filePath,
    content: fs.readFileSync(filePath, 'utf-8'),
    testPath: filePath.replace('.tsx', '.test.tsx')
  };
}

// Write component file
function writeComponent(name, content) {
  const fileName = name.toLowerCase() + '.tsx';
  const filePath = path.join(UI_DIR, fileName);
  
  if (!fs.existsSync(UI_DIR)) {
    fs.mkdirSync(UI_DIR, { recursive: true });
  }
  
  fs.writeFileSync(filePath, content);
  return filePath;
}

// Write test file
function writeTest(name) {
  const fileName = name.toLowerCase() + '.test.tsx';
  const filePath = path.join(UI_DIR, fileName);
  const content = TEST_TEMPLATE(name);
  
  fs.writeFileSync(filePath, content);
  return filePath;
}

// Main command handler
function main() {
  // Get input from command line or stdin
  const input = process.argv.slice(2).join(' ') || 
                require('fs').readFileSync(0, 'utf-8').trim();
  
  if (!input) {
    console.log(`Usage: /component <description>

Examples:
  /component create a card for displaying server info
  /component update button to add loading state
  /component fix combobox missing optionLabel prop
  /component make a dropdown with search functionality`);
    process.exit(1);
  }
  
  const request = parseRequest(input);
  
  if (request.type === 'update') {
    const existing = readComponent(request.componentName);
    if (!existing) {
      console.error(`❌ Component "${request.componentName}" not found in ${UI_DIR}`);
      console.error(`   Available components:`);
      const files = fs.readdirSync(UI_DIR).filter(f => f.endsWith('.tsx') && !f.startsWith('_'));
      files.forEach(f => console.error(`   - ${f.replace('.tsx', '')}`));
      process.exit(1);
    }
    
    console.log(`✏️  Update request for ${request.componentName}`);
    console.log(`   Description: ${request.description}`);
    console.log(`   File: ${existing.path}`);
    console.log(`\n⚠️  Manual update required - skill guidance loaded from:`);
    console.log(`   ${SKILL_PATH}`);
    console.log(`\n📋 Apply these conventions:`);
    console.log(`   - Use Kobalte primitives for interactive elements`);
    console.log(`   - Style with Tailwind using theme variables`);
    console.log(`   - Follow existing component patterns`);
    console.log(`   - Update tests to cover new functionality`);
    
  } else {
    // Create new component
    const componentName = request.componentName || generateName(request.description);
    const pattern = detectPattern(request.description);
    
    // Check if component already exists
    const existing = readComponent(componentName);
    if (existing) {
      console.error(`❌ Component "${componentName}" already exists at ${existing.path}`);
      console.error(`   Use "/component update ${componentName} ..." to modify it`);
      process.exit(1);
    }
    
    // Determine Kobalte primitive for interactive components
    let primitive = 'Button';
    if (pattern === 'kobalte') {
      const normalized = request.description.toLowerCase();
      if (normalized.includes('select')) primitive = 'Select';
      else if (normalized.includes('dialog') || normalized.includes('modal')) primitive = 'Dialog';
      else if (normalized.includes('tooltip')) primitive = 'Tooltip';
      else if (normalized.includes('popover')) primitive = 'Popover';
      else if (normalized.includes('checkbox')) primitive = 'Checkbox';
      else if (normalized.includes('radio')) primitive = 'RadioGroup';
      else if (normalized.includes('slider')) primitive = 'Slider';
      else if (normalized.includes('switch')) primitive = 'Switch';
      else if (normalized.includes('tabs')) primitive = 'Tabs';
    }
    
    // Generate component code
    let code;
    if (pattern === 'simple') {
      code = PATTERNS.simple(componentName, request.description);
    } else if (pattern === 'kobalte') {
      code = PATTERNS.kobalte(componentName, primitive, request.description);
    } else {
      code = PATTERNS.cva(componentName, request.description);
    }
    
    // Write files
    const componentPath = writeComponent(componentName, code);
    const testPath = writeTest(componentName);
    
    console.log(`✅ Created component: ${componentName}`);
    console.log(`   Pattern: ${pattern}`);
    console.log(`   Component: ${componentPath}`);
    console.log(`   Test: ${testPath}`);
    console.log(`\n📋 Usage:`);
    console.log(`   import { ${componentName} } from "~/components/ui/${componentName.toLowerCase()}"`);
    console.log(`\n⚠️  Note: Test framework (vitest) not yet configured in this project.`);
    console.log(`   Tests generated but won't run until vitest is set up.`);
  }
}

if (require.main === module) {
  main();
}

module.exports = { parseRequest, detectPattern, generateName, toPascalCase };
