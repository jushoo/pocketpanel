import { query, createAsync, redirect, useParams } from "@solidjs/router";
import { Show } from "solid-js";
import { Button } from "~/components/ui/button";
import { Badge } from "~/components/ui/badge";
import { ModeToggle } from "~/components/ModeToggle";
import { ArrowLeft, Circle } from "lucide-solid";
import { formatMemory } from "~/lib/utils";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3001";

interface ServerData {
  id: number;
  name: string;
  type: string;
  version: string;
  min_mem: number;
  max_mem: number;
  port: number;
}

// Query to fetch a single server
const getServer = query(async (id: string) => {
  const response = await fetch(`${API_URL}/api/v1/servers/${id}`, {
    credentials: "include",
  });
  
  if (!response.ok) {
    if (response.status === 401) {
      throw redirect("/");
    }
    if (response.status === 404) {
      throw new Error("Server not found");
    }
    throw new Error("Failed to fetch server");
  }
  
  return response.json() as Promise<ServerData>;
}, "getServer");

// Route data loader - must be exported
export function routeData() {
  const params = useParams();
  return {
    server: createAsync(() => getServer(params.id)),
  };
}

export default function ServerDetailPage() {
  // Access route data
  const params = useParams();
  const server = createAsync(() => getServer(params.id));

  return (
    <div class="min-h-screen bg-background">
      <header class="border-b">
        <div class="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
          <div class="flex items-center gap-4">
            <Button variant="ghost" size="sm" as="a" href="/dashboard">
              <ArrowLeft class="mr-2 h-4 w-4" />
              Back to servers
            </Button>
          </div>
          <div class="flex items-center gap-4">
            <ModeToggle />
          </div>
        </div>
      </header>

      <main class="mx-auto max-w-4xl px-6 py-8">
        <Show when={!server.loading} fallback={<div class="text-center py-12">Loading...</div>}>
          <Show when={!server.error} fallback={<div class="text-center py-12 text-destructive">Failed to load server</div>}>
            <Show when={server()} fallback={<div class="text-center py-12">Server not found</div>}>
              {(s) => (
                <>
                  {/* Title */}
                  <div class="mb-8">
                    <h1 class="text-xl font-medium text-foreground">{s().name}</h1>
                    <p class="mt-1 text-sm text-muted-foreground">Server configuration and details</p>
                  </div>

                  {/* Server Info - List Style */}
                  <div class="space-y-6">
                    {/* Name */}
                    <div class="pb-6 border-b border-border">
                      <div class="mb-2 flex items-center justify-between">
                        <span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Name</span>
                      </div>
                      <p class="text-lg font-medium text-foreground">{s().name}</p>
                    </div>

                    {/* Type & Version Row */}
                    <div class="pb-6 border-b border-border">
                      <span class="mb-2 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Version</span>
                      <div class="flex items-center gap-3">
                        <p class="text-foreground">v{s().version}</p>
                        <Badge variant="secondary" class="text-xs">{s().type}</Badge>
                      </div>
                    </div>

                    {/* Resources */}
                    <div class="pb-6 border-b border-border">
                      <span class="mb-3 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Resources</span>
                      <div class="grid grid-cols-2 gap-6">
                        <div>
                          <p class="text-xs text-muted-foreground">Memory</p>
                          <p class="mt-1 font-medium text-foreground">
                            {formatMemory(s().min_mem)} - {formatMemory(s().max_mem)}
                          </p>
                        </div>
                        <div>
                          <p class="text-xs text-muted-foreground">Port</p>
                          <p class="mt-1 font-medium text-foreground">{s().port}</p>
                        </div>
                      </div>
                    </div>

                    {/* Status */}
                    <div class="pb-6 border-b border-border">
                      <span class="mb-2 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Status</span>
                      <div class="flex items-center gap-2">
                        <Circle class="h-3 w-3 text-yellow-500" fill="currentColor" />
                        <span class="text-sm text-muted-foreground">Unknown</span>
                      </div>
                    </div>
                  </div>
                </>
              )}
            </Show>
          </Show>
        </Show>
      </main>
    </div>
  );
}
