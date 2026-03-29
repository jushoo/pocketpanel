import { useColorMode } from "@kobalte/core";
import { Sun, Moon } from "lucide-solid";
import { Button } from "~/components/ui/button";

export function ModeToggle() {
  const { colorMode, setColorMode } = useColorMode();

  const toggleMode = () => {
    setColorMode(colorMode() === "light" ? "dark" : "light");
  };

  return (
    <div class="fixed top-4 right-4 z-50">
      <Button onClick={toggleMode} variant="outline" size="icon">
        <Sun
          class={`h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all ${colorMode() === "dark" ? "scale-0 -rotate-90" : ""}`}
        />
        <Moon
          class={`absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all ${colorMode() === "dark" ? "scale-100 rotate-0" : ""}`}
        />
        <span class="sr-only">Toggle theme</span>
      </Button>
    </div>
  );
}
