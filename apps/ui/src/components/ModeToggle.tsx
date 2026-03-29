import * as DropdownMenu from "@kobalte/core/dropdown-menu";

import { Moon, Sun, Laptop } from "lucide-solid";

import { Button } from "~/components/ui/button";
import { useColorMode } from "~/lib/color-mode";
import { cn } from "~/lib/utils";

export function ModeToggle() {
  const { colorMode, setColorMode } = useColorMode();

  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        as={(props) => (
          <Button variant="ghost" size="sm" class="w-9 px-0" {...props}>
            <Sun class="size-6 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
            <Moon class="absolute size-6 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
            <span class="sr-only">Toggle theme</span>
          </Button>
        )}
      />
      <DropdownMenu.Portal>
        <DropdownMenu.Content class="z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md">
          <DropdownMenu.Item
            class={cn(
              "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
              colorMode() === "light" && "bg-accent text-accent-foreground"
            )}
            onSelect={() => setColorMode("light")}
          >
            <Sun class="mr-2 size-4" />
            <span>Light</span>
          </DropdownMenu.Item>
          <DropdownMenu.Item
            class={cn(
              "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
              colorMode() === "dark" && "bg-accent text-accent-foreground"
            )}
            onSelect={() => setColorMode("dark")}
          >
            <Moon class="mr-2 size-4" />
            <span>Dark</span>
          </DropdownMenu.Item>
          <DropdownMenu.Item
            class={cn(
              "relative flex cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground focus:bg-accent focus:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50",
              colorMode() === "system" && "bg-accent text-accent-foreground"
            )}
            onSelect={() => setColorMode("system")}
          >
            <Laptop class="mr-2 size-4" />
            <span>System</span>
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
}
