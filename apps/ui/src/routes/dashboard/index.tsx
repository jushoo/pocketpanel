import { query, createAsync, useSubmission, action, redirect, useParams } from "@solidjs/router";
import { createSignal, Show, For } from "solid-js";
import { Button } from "~/components/ui/button";
import { Badge } from "~/components/ui/badge";
import { ModeToggle } from "~/components/ModeToggle";
import { Plus, Server as ServerIcon, Circle } from "lucide-solid";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3001";

interface ServerData {
  id: number;
  name: string;
  type: string;
  version: string;
  min_mem: number;
  max_mem: number;
  port: number;
  status?: string;
}

// Query to fetch servers
const getServers = query(async () => {
  const response = await fetch(`${API_URL}/api/v1/servers`, {
    credentials: "include",
  });
  
  if (!response.ok) {
    if (response.status === 401) {
      throw redirect("/");
    }
    throw new Error("Failed to fetch servers");
  }
  
  return response.json() as Promise<ServerData[]>;
}, "getServers");

// Query to fetch user
const getUser = query(async () => {
  const response = await fetch(`${API_URL}/api/v1/me`, {
    credentials: "include",
  });
  
  if (!response.ok) {
    return null;
  }
  
  return response.json();
}, "getUser");

// Action to logout
const logoutAction = action(async () => {
  try {
    await fetch(`${API_URL}/api/v1/auth/logout`, {
      method: "POST",
      credentials: "include",
    });
  } catch {
    // Continue even if logout fails
  }
  
  throw redirect("/");
}, "logout");

// Route data loader - must be exported
export function routeData() {
  return {
    servers: createAsync(() => getServers()),
    user: createAsync(() => getUser()),
  };
}

export default function DashboardPage() {
  // Access route data - in SolidStart we need to access it differently
  // For now, let's use createAsync directly in the component
  const servers = createAsync(() => getServers());
  const user = createAsync(() => getUser());
  
  const logoutSubmission = useSubmission(logoutAction);
  const [activeFilter, setActiveFilter] = createSignal<"all" | "vanilla" | "fabric">("all");

  const filteredServers = () => {
    const serverList = servers() || [];
    if (activeFilter() === "all") {
      return serverList;
    }
    return serverList.filter(
      (s: ServerData) => s.type.toLowerCase() === activeFilter(),
    );
  };

  const getStatusColor = (status?: string) => {
    switch (status?.toLowerCase()) {
      case "running":
        return "text-green-500";
      case "stopped":
        return "text-muted-foreground";
      default:
        return "text-yellow-500";
    }
  };

  return (
    <div class="min-h-screen bg-background">
      <header class="border-b">
        <div class="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
          <div>
            <h1 class="text-xl font-medium text-foreground">Servers</h1>
            <p class="text-sm text-muted-foreground">
              Manage your Minecraft servers
            </p>
          </div>
          <div class="flex items-center gap-4">
            <ModeToggle />
            <Show when={user()}>
              <span class="text-sm text-muted-foreground">
                {user()?.username}
              </span>
            </Show>
            <form action={logoutAction} method="post">
              <Button type="submit" variant="ghost" size="sm" disabled={logoutSubmission.pending}>
                {logoutSubmission.pending ? "Logging out..." : "Log out"}
              </Button>
            </form>
          </div>
        </div>
      </header>

      <main class="mx-auto max-w-4xl px-6 py-8">
        <Show
          when={!servers.loading}
          fallback={<div class="text-center py-12">Loading...</div>}
        >
          <Show
            when={!servers.error}
            fallback={
              <div class="flex min-h-[60vh] flex-col items-center justify-center text-center">
                <p class="text-sm text-destructive">Failed to load servers</p>
                <Button 
                  class="mt-4" 
                  variant="outline" 
                  onClick={() => window.location.reload()}
                >
                  Retry
                </Button>
              </div>
            }
          >
            <Show
              when={filteredServers().length > 0}
              fallback={
                <div class="flex min-h-[60vh] flex-col items-center justify-center text-center">
                  <div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-muted">
                    <ServerIcon class="h-6 w-6 text-muted-foreground" />
                  </div>
                  <h2 class="text-lg font-medium text-foreground">
                    No servers yet
                  </h2>
                  <p class="mt-1 text-sm text-muted-foreground">
                    Create your first Minecraft server to get started
                  </p>
                  <Button class="mt-6" as="a" href="/dashboard/servers/new">
                    <Plus class="mr-2 h-4 w-4" />
                    Create Server
                  </Button>
                </div>
              }
            >
              <div class="mb-6 flex items-center justify-between">
                <p class="text-sm text-muted-foreground">
                  {filteredServers().length} server
                  {filteredServers().length === 1 ? "" : "s"}
                </p>
                <Button size="sm" as="a" href="/dashboard/servers/new">
                  <Plus class="mr-2 h-4 w-4" />
                  New Server
                </Button>
              </div>

              {/* Filter Tags */}
              <div class="mb-8 flex items-center gap-2">
                <button
                  type="button"
                  onClick={() => setActiveFilter("all")}
                  class={`rounded-full border px-4 py-1.5 text-sm transition-colors ${
                    activeFilter() === "all"
                      ? "border-primary bg-primary text-primary-foreground"
                      : "border-input bg-background text-foreground hover:bg-muted"
                  }`}
                >
                  All
                </button>
                <button
                  type="button"
                  onClick={() => setActiveFilter("vanilla")}
                  class={`rounded-full border px-4 py-1.5 text-sm transition-colors ${
                    activeFilter() === "vanilla"
                      ? "border-primary bg-primary text-primary-foreground"
                      : "border-input bg-background text-foreground hover:bg-muted"
                  }`}
                >
                  Vanilla
                </button>
                <button
                  type="button"
                  onClick={() => setActiveFilter("fabric")}
                  class={`rounded-full border px-4 py-1.5 text-sm transition-colors ${
                    activeFilter() === "fabric"
                      ? "border-primary bg-primary text-primary-foreground"
                      : "border-input bg-background text-foreground hover:bg-muted"
                  }`}
                >
                  Fabric
                </button>
              </div>

              {/* Server List */}
              <div class="space-y-6">
                <For each={filteredServers()}>
                  {(server) => (
                    <a
                      href={`/dashboard/servers/${server.id}`}
                      class="group block border-b border-border pb-6 transition-colors last:border-0"
                    >
                      <div class="space-y-3">
                        <h3 class="text-xl font-medium text-foreground">
                          {server.name}
                        </h3>
                        <p class="text-sm leading-relaxed text-muted-foreground">
                          Minecraft server running on port {server.port} with{" "}
                          {server.min_mem}MB - {server.max_mem}MB memory
                          allocation
                        </p>
                        <div class="flex items-center gap-2 pt-1">
                          <Badge variant="secondary" class="text-xs">
                            {server.type}
                          </Badge>
                          <Badge variant="outline" class="text-xs">
                            v{server.version}
                          </Badge>
                        </div>
                        <div class="flex items-center gap-4 pt-2 text-xs text-muted-foreground">
                          <div class="flex items-center gap-1.5">
                            <Circle
                              class={`h-3 w-3 ${getStatusColor(server.status)}`}
                              fill="currentColor"
                            />
                            <span class="capitalize">
                              {server.status || "Unknown"}
                            </span>
                          </div>
                          <span>Port {server.port}</span>
                        </div>
                      </div>
                    </a>
                  )}
                </For>
              </div>
            </Show>
          </Show>
        </Show>
      </main>
    </div>
  );
}
