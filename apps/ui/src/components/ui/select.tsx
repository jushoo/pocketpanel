import { splitProps, type JSX } from "solid-js"
import { Select as SelectPrimitive } from "@kobalte/core/select"
import { cn } from "~/lib/utils"

const Select = SelectPrimitive

const SelectTrigger = (
  props: JSX.HTMLAttributes<HTMLButtonElement> & { class?: string }
) => {
  const [local, others] = splitProps(props, ["class", "children"])
  return (
    <SelectPrimitive.Trigger
      class={cn(
        "flex h-11 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1",
        local.class
      )}
      {...others}
    >
      {local.children}
      <SelectPrimitive.Icon class="opacity-50">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          width="24"
          height="24"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          stroke-width="2"
          stroke-linecap="round"
          stroke-linejoin="round"
          class="h-4 w-4"
        >
          <path d="m6 9 6 6 6-6" />
        </svg>
      </SelectPrimitive.Icon>
    </SelectPrimitive.Trigger>
  )
}

const SelectContent = (
  props: JSX.HTMLAttributes<HTMLDivElement> & { class?: string }
) => {
  const [local, others] = splitProps(props, ["class"])
  return (
    <SelectPrimitive.Portal>
      <SelectPrimitive.Content
        class={cn(
          "relative z-50 max-h-[--kb-select-content-available-height] min-w-[8rem] overflow-y-auto overflow-x-hidden rounded-md border bg-popover text-popover-foreground shadow-md data-[expanded]:animate-in data-[closed]:animate-out data-[closed]:fade-out-0 data-[expanded]:fade-in-0 data-[closed]:zoom-out-95 data-[expanded]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 origin-[--kb-select-content-transform-origin]",
          local.class
        )}
        {...others}
      >
        <SelectPrimitive.Listbox class="p-1" />
      </SelectPrimitive.Content>
    </SelectPrimitive.Portal>
  )
}

const SelectItem = (
  props: { item: any; class?: string; children?: any }
) => {
  const [local, others] = splitProps(props, ["class", "children", "item"])
  return (
    <SelectPrimitive.Item
      item={local.item}
      class={cn(
        "relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
        local.class
      )}
      {...others}
    >
      <span class="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
        <SelectPrimitive.ItemIndicator>
          <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
            class="h-4 w-4"
          >
            <path d="M20 6 9 17l-5-5" />
          </svg>
        </SelectPrimitive.ItemIndicator>
      </span>
      <SelectPrimitive.ItemLabel>{local.children}</SelectPrimitive.ItemLabel>
    </SelectPrimitive.Item>
  )
}

export { Select, SelectTrigger, SelectContent, SelectItem }
