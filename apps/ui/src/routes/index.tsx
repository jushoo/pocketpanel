import { action, useSubmission, redirect } from "@solidjs/router";
import { Show } from "solid-js";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import { ModeToggle } from "~/components/ModeToggle";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3001";

const loginAction = action(async (formData: FormData) => {
  const username = formData.get("username") as string;
  const password = formData.get("password") as string;

  if (!username || !password) {
    return { ok: false, message: "Please enter both username and password" };
  }

  try {
    const response = await fetch(`${API_URL}/api/v1/auth/login`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ username, password }),
    });

    if (!response.ok) {
      const data = await response.json();
      return { ok: false, message: data.error || "Login failed" };
    }

    throw redirect("/dashboard");
  } catch (err) {
    if (err instanceof Response && err.status === 302) {
      throw err;
    }
    return { ok: false, message: "An error occurred during login" };
  }
}, "login");

export default function LoginPage() {
  const submission = useSubmission(loginAction);

  return (
    <div class="relative flex min-h-screen w-full items-center justify-center bg-background">
      <div class="absolute right-4 top-4">
        <ModeToggle />
      </div>
      <div class="w-full max-w-sm px-6">
        <div class="mb-12 text-center">
          <div class="mb-6 inline-flex h-12 w-12 items-center justify-center rounded-xl border border-border bg-muted">
            <svg
              viewBox="0 0 24 24"
              class="h-6 w-6 text-foreground"
              fill="none"
              stroke="currentColor"
              stroke-width="1.5"
            >
              <rect x="3" y="3" width="18" height="18" rx="2" />
              <path d="M9 12h6M12 9v6" />
            </svg>
          </div>
          <h1 class="text-xl font-medium tracking-tight text-foreground">
            PocketPanel
          </h1>
          <p class="mt-2 text-sm text-muted-foreground">Admin Login</p>
        </div>

        <form action={loginAction} method="post" class="space-y-5">
          <div class="space-y-2">
            <Label
              for="username"
              class="text-sm font-normal text-muted-foreground"
            >
              Username
            </Label>
            <Input
              id="username"
              name="username"
              type="text"
              class="h-11"
              required
            />
          </div>

          <div class="space-y-2">
            <Label
              for="password"
              class="text-sm font-normal text-muted-foreground"
            >
              Password
            </Label>
            <Input
              id="password"
              name="password"
              type="password"
              class="h-11"
              required
            />
          </div>

          <Show when={submission.result?.ok === false}>
            <p class="text-sm text-destructive">{submission.result?.message}</p>
          </Show>

          <Button type="submit" disabled={submission.pending} class="h-11 w-full">
            {submission.pending ? "Signing in..." : "Sign in"}
          </Button>
        </form>

        <p class="mt-8 text-center text-xs text-muted-foreground">
          Secure admin access only
        </p>
      </div>
    </div>
  );
}
