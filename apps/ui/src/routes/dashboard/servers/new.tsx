import {
  query,
  createAsync,
  action,
  redirect,
  useSubmission,
} from "@solidjs/router";
import { createSignal, Show } from "solid-js";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Label } from "~/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
import { Boxes, Info, ArrowLeft } from "lucide-solid";
import {
  Combobox,
  ComboboxContent,
  ComboboxControl,
  ComboboxHiddenSelect,
  ComboboxInput,
  ComboboxItem,
  ComboboxItemIndicator,
  ComboboxItemLabel,
  ComboboxTrigger,
} from "~/components/ui/combobox";

const API_URL = import.meta.env.VITE_API_URL || "http://localhost:3001";

interface ServerType {
  id: string;
  name: string;
  description: string;
}

// Query to fetch server types
const getServerTypes = query(async () => {
  const response = await fetch(`${API_URL}/api/v1/server-types`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch server types");
  }

  return response.json() as Promise<ServerType[]>;
}, "getServerTypes");

// Query to fetch versions for a type
const getVersions = query(async (typeId: string) => {
  if (!typeId) return [];

  const response = await fetch(`${API_URL}/api/v1/versions/${typeId}`, {
    credentials: "include",
  });

  if (!response.ok) {
    throw new Error("Failed to fetch versions");
  }

  const data = await response.json();
  return data.versions as string[];
}, "getVersions");

// Action to create server
const createServerAction = action(async (formData: FormData) => {
  const name = formData.get("name") as string;
  const serverType = formData.get("serverType") as string;
  const version = formData.get("version") as string;
  const port = formData.get("port") as string;
  const minMemory = formData.get("minMemory") as string;
  const maxMemory = formData.get("maxMemory") as string;

  try {
    const response = await fetch(`${API_URL}/api/v1/servers`, {
      method: "POST",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        name,
        type: serverType,
        version,
        port: port ? parseInt(port) : undefined,
        min_mem: parseInt(minMemory || "2"),
        max_mem: parseInt(maxMemory || "4"),
      }),
    });

    if (!response.ok) {
      const data = await response.json();
      return { ok: false, message: data.error || "Failed to create server" };
    }

    const result = await response.json();
    throw redirect(`/dashboard/servers/${result.id}`);
  } catch (err) {
    if (err instanceof Response && err.status === 302) {
      throw err;
    }
    return {
      ok: false,
      message: "An error occurred while creating the server",
    };
  }
}, "createServer");

// Route data loader - must be exported
export function routeData() {
  return {
    serverTypes: createAsync(() => getServerTypes()),
  };
}

export default function CreateServerPage() {
  // Access route data via createAsync (following codebase pattern)
  const serverTypes = createAsync(() => getServerTypes());

  const submission = useSubmission(createServerAction);
  const [selectedTypeId, setSelectedTypeId] = createSignal("");
  const [selectedVersion, setSelectedVersion] = createSignal("");
  const [nameLength, setNameLength] = createSignal(0);

  const versions = createAsync(() => getVersions(selectedTypeId()));

  const selectedType = () =>
    serverTypes()?.find((t: ServerType) => t.id === selectedTypeId());

  const handleNameInput = (e: InputEvent) => {
    const target = e.target as HTMLInputElement;
    setNameLength(target.value.length);
  };

  return (
    <div class="min-h-screen bg-background">
      <header class="border-b">
        <div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <div class="flex items-center gap-4">
            <Button variant="ghost" size="sm" as="a" href="/dashboard">
              <ArrowLeft class="mr-2 h-4 w-4" />
              Back
            </Button>
            <div>
              <h1 class="text-xl font-medium text-foreground">Create Server</h1>
              <p class="text-sm text-muted-foreground">
                Set up a new Minecraft server
              </p>
            </div>
          </div>
        </div>
      </header>

      <main class="mx-auto max-w-6xl px-6 py-8">
        <div class="flex min-h-[60vh] items-center justify-center">
          <Card class="w-full max-w-md">
            <CardHeader class="text-center">
              <div class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-muted">
                <Boxes class="h-6 w-6 text-muted-foreground" />
              </div>
              <CardTitle class="text-lg font-medium">Create Server</CardTitle>
              <CardDescription>
                Configure your new Minecraft server
              </CardDescription>
            </CardHeader>
            <CardContent>
              {submission.result?.ok === false && (
                <div class="mb-4 rounded-lg border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
                  {submission.result.message}
                </div>
              )}

              <form action={createServerAction} method="post" class="space-y-6">
                {/* Server Name */}
                <div class="space-y-2">
                  <Label for="name">
                    Server Name
                    <span class="text-destructive">*</span>
                  </Label>
                  <Input
                    id="name"
                    name="name"
                    type="text"
                    placeholder="My Minecraft Server"
                    maxLength={100}
                    onInput={handleNameInput}
                    class="h-11"
                    required
                  />
                  <div class="flex justify-end">
                    <span class="text-xs text-muted-foreground">
                      {nameLength()} / 100
                    </span>
                  </div>
                </div>

                {/* Server Type */}
                <div class="space-y-2">
                  <Label for="serverType">
                    Server Type
                    <span class="text-destructive">*</span>
                  </Label>
                  <Show
                    when={!serverTypes.loading}
                    fallback={<div>Loading types...</div>}
                  >
                    <Combobox
                      options={serverTypes() || []}
                      optionValue="id"
                      optionTextValue="name"
                      placeholder="Choose a server type"
                      value={serverTypes()?.find((t: ServerType) => t.id === selectedTypeId()) || null}
                      onChange={(selected) => {
                        setSelectedTypeId(selected?.id || "");
                        setSelectedVersion("");
                      }}
                      itemComponent={(props) => (
                        <ComboboxItem item={props.item}>
                          <ComboboxItemLabel>{props.item.rawValue.name}</ComboboxItemLabel>
                          <ComboboxItemIndicator />
                        </ComboboxItem>
                      )}
                    >
                      <ComboboxHiddenSelect name="serverType" required />
                      <ComboboxControl class="h-11 w-full">
                        <ComboboxInput />
                        <ComboboxTrigger />
                      </ComboboxControl>
                      <ComboboxContent />
                    </Combobox>
                  </Show>
                  <Show when={selectedType()}>
                    <p class="text-xs text-muted-foreground">
                      {selectedType()?.description}
                    </p>
                  </Show>
                </div>

                {/* Server Version */}
                <div class="space-y-2">
                  <Label for="version">
                    Server Version
                    <span class="text-destructive">*</span>
                  </Label>
                  <Show
                    when={!versions.loading}
                    fallback={<div>Loading versions...</div>}
                  >
                    <Combobox
                      options={versions() || []}
                      placeholder={selectedTypeId() ? "Choose a version" : "Select a server type first"}
                      value={selectedVersion() || null}
                      onChange={(selected) => setSelectedVersion(selected || "")}
                      disabled={!selectedTypeId()}
                      itemComponent={(props) => (
                        <ComboboxItem item={props.item}>
                          <ComboboxItemLabel>{props.item.rawValue}</ComboboxItemLabel>
                          <ComboboxItemIndicator />
                        </ComboboxItem>
                      )}
                    >
                      <ComboboxHiddenSelect name="version" required />
                      <ComboboxControl class="h-11 w-full">
                        <ComboboxInput />
                        <ComboboxTrigger />
                      </ComboboxControl>
                      <ComboboxContent />
                    </Combobox>
                  </Show>
                </div>

                {/* Divider */}
                <div class="border-t"></div>

                {/* Port */}
                <div class="space-y-2">
                  <Label for="port">Port</Label>
                  <Input
                    id="port"
                    name="port"
                    type="number"
                    placeholder="25565"
                    min={1024}
                    max={65535}
                    class="h-11"
                  />
                  <p class="text-xs text-muted-foreground">
                    Auto-assigned if left empty
                  </p>
                </div>

                {/* Memory Allocation */}
                <div class="space-y-3">
                  <Label>Memory Allocation (GB)</Label>
                  <div class="grid grid-cols-2 gap-4">
                    <div class="space-y-2">
                      <Label
                        for="minMemory"
                        class="text-xs text-muted-foreground"
                      >
                        Minimum
                      </Label>
                      <Input
                        id="minMemory"
                        name="minMemory"
                        type="number"
                        placeholder="2"
                        min={1}
                        defaultValue="2"
                        class="h-11"
                      />
                    </div>
                    <div class="space-y-2">
                      <Label
                        for="maxMemory"
                        class="text-xs text-muted-foreground"
                      >
                        Maximum
                      </Label>
                      <Input
                        id="maxMemory"
                        name="maxMemory"
                        type="number"
                        placeholder="4"
                        min={1}
                        max={128}
                        defaultValue="4"
                        class="h-11"
                      />
                    </div>
                  </div>
                </div>

                {/* Info Note */}
                <div class="flex items-start gap-2 rounded-lg bg-muted px-3 py-2">
                  <Info class="mt-0.5 h-4 w-4 flex-shrink-0 text-muted-foreground" />
                  <p class="text-xs text-muted-foreground">
                    These settings can be changed later
                  </p>
                </div>

                {/* Actions */}
                <div class="flex items-center justify-between gap-4 pt-2">
                  <Button
                    type="button"
                    variant="ghost"
                    as="a"
                    href="/dashboard"
                    disabled={submission.pending}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    disabled={submission.pending}
                    class="min-w-[140px]"
                  >
                    {submission.pending ? "Creating..." : "Create Server"}
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </main>
    </div>
  );
}
