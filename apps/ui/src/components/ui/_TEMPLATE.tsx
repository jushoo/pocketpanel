/**
 * COMPONENT TEMPLATE
 * 
 * This file serves as a template for creating new UI components.
 * Copy this file and adapt it to your needs using one of the three patterns below.
 */

import type { JSX, ValidComponent, Component, ComponentProps } from "solid-js"
import { splitProps } from "solid-js"
import type { VariantProps } from "class-variance-authority"
import { cva } from "class-variance-authority"
import type { PolymorphicProps } from "@kobalte/core/polymorphic"

import { cn } from "~/lib/utils"

// ============================================================================
// PATTERN 1: Simple HTML Wrapper
// For basic styled elements (Input, Label, Badge, Card)
// ============================================================================

const SimpleExample: Component<ComponentProps<"div">> = (props) => {
  const [local, others] = splitProps(props, ["class"])
  return (
    <div
      class={cn(
        "rounded-md border bg-background p-4",
        local.class
      )}
      {...others}
    />
  )
}

export { SimpleExample }

// ============================================================================
// PATTERN 2: Kobalte Primitive Wrapper
// For complex interactive components (Button, Select, Combobox)
// ============================================================================

// Uncomment and adapt for Kobalte-based components:
// import * as ButtonPrimitive from "@kobalte/core/button"
//
// type KobalteExampleProps<T extends ValidComponent = "button"> = 
//   ButtonPrimitive.ButtonRootProps<T> & { class?: string }
//
// const KobalteExample = <T extends ValidComponent = "button">(
//   props: PolymorphicProps<T, KobalteExampleProps<T>>
// ) => {
//   const [local, others] = splitProps(props as KobalteExampleProps, ["class"])
//   return (
//     <ButtonPrimitive.Root
//       class={cn("inline-flex items-center...", local.class)}
//       {...others}
//     />
//   )
// }
//
// export { KobalteExample }

// ============================================================================
// PATTERN 3: CVA Variants
// For components with multiple visual states (Button, Badge)
// ============================================================================

// Uncomment and adapt for variant-based components:
// const exampleVariants = cva(
//   "inline-flex items-center justify-center...",
//   {
//     variants: {
//       variant: {
//         default: "bg-primary text-primary-foreground",
//         outline: "border border-input"
//       },
//       size: {
//         default: "h-10 px-4 py-2",
//         sm: "h-9 px-3"
//       }
//     },
//     defaultVariants: {
//       variant: "default",
//       size: "default"
//     }
//   }
// )
//
// type CVAExampleProps = ComponentProps<"div"> &
//   VariantProps<typeof exampleVariants>
//
// const CVAExample: Component<CVAExampleProps> = (props) => {
//   const [local, others] = splitProps(props, ["class", "variant", "size"])
//   return (
//     <div
//       class={cn(exampleVariants({ variant: local.variant, size: local.size }), local.class)}
//       {...others}
//     />
//   )
// }
//
// export { CVAExample, exampleVariants }
