import js from "@eslint/js";
import tsParser from "@typescript-eslint/parser";
import tsPlugin from "@typescript-eslint/eslint-plugin";
import solid from "eslint-plugin-solid";
import globals from "globals";

export default [
  js.configs.recommended,
  {
    ignores: [
      "**/node_modules/**",
      "**/.output/**",
      "**/.nitro/**",
      "**/dist/**",
      "**/build/**",
      "**/*.d.ts",
    ],
  },
  {
    files: ["**/*.{ts,tsx}"],
    languageOptions: {
      parser: tsParser,
      parserOptions: {
        project: "./tsconfig.json",
        tsconfigRootDir: import.meta.dirname,
      },
      globals: {
        ...globals.browser,
        ...globals.node,
      },
    },
    plugins: {
      "@typescript-eslint": tsPlugin,
      solid: solid,
    },
    rules: {
      ...tsPlugin.configs.recommended.rules,
      ...solid.configs.recommended.rules,
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
      "solid/no-destructure": "error",
    },
  },
  {
    files: ["**/*.config.{ts,js}"],
    languageOptions: {
      parser: tsParser,
    },
    rules: {
      "@typescript-eslint/no-require-imports": "off",
    },
  },
];
