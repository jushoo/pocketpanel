import { createContext, useContext, createSignal, createEffect, type JSX } from "solid-js";

type ColorMode = "light" | "dark" | "system";

interface ColorModeContextType {
  colorMode: () => ColorMode;
  setColorMode: (mode: ColorMode) => void;
  resolvedMode: () => "light" | "dark";
}

const ColorModeContext = createContext<ColorModeContextType>();

const STORAGE_KEY = "color-mode";

function getStoredMode(): ColorMode {
  if (typeof document === "undefined") return "system";
  const cookie = document.cookie.match(/color-mode=([^;]+)/);
  if (cookie) return cookie[1] as ColorMode;
  const local = localStorage.getItem(STORAGE_KEY);
  if (local) return local as ColorMode;
  return "system";
}

function setStoredMode(mode: ColorMode) {
  if (typeof document === "undefined") return;
  // Set cookie for SSR
  document.cookie = `${STORAGE_KEY}=${mode};path=/;max-age=31536000`;
  // Set localStorage for persistence
  localStorage.setItem(STORAGE_KEY, mode);
}

function getResolvedMode(mode: ColorMode): "light" | "dark" {
  if (mode === "dark") return "dark";
  if (mode === "light") return "light";
  // system mode
  if (typeof window !== "undefined") {
    return window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
  }
  return "light";
}

export function ColorModeProvider(props: { children: JSX.Element }) {
  const [colorMode, setMode] = createSignal<ColorMode>(getStoredMode());
  const [resolvedMode, setResolvedMode] = createSignal<"light" | "dark">(
    getResolvedMode(getStoredMode())
  );

  const setColorMode = (mode: ColorMode) => {
    setMode(mode);
    setStoredMode(mode);
    setResolvedMode(getResolvedMode(mode));
  };

  createEffect(() => {
    const resolved = resolvedMode();
    
    // Apply class to html element
    const html = document.documentElement;
    if (resolved === "dark") {
      html.classList.add("dark");
    } else {
      html.classList.remove("dark");
    }
  });

  // Listen for system changes
  createEffect(() => {
    if (colorMode() !== "system") return;
    
    const mediaQuery = window.matchMedia("(prefers-color-scheme: dark)");
    const handler = (e: MediaQueryListEvent) => {
      setResolvedMode(e.matches ? "dark" : "light");
    };
    
    mediaQuery.addEventListener("change", handler);
    return () => mediaQuery.removeEventListener("change", handler);
  });

  return (
    <ColorModeContext.Provider value={{ colorMode, setColorMode, resolvedMode }}>
      {props.children}
    </ColorModeContext.Provider>
  );
}

export function useColorMode() {
  const context = useContext(ColorModeContext);
  if (!context) {
    throw new Error("useColorMode must be used within a ColorModeProvider");
  }
  return context;
}

// Script to run before hydration
export const colorModeScript = `
  (function() {
    try {
      const mode = document.cookie.match(/color-mode=([^;]+)/)?.[1] || localStorage.getItem("color-mode") || "system";
      let resolved = mode;
      if (mode === "system") {
        resolved = window.matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
      }
      if (resolved === "dark") {
        document.documentElement.classList.add("dark");
      }
    } catch (e) {}
  })();
`;
