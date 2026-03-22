<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { Plus, Server as ServerIcon, Circle } from '@lucide/svelte';
	import { enhance } from '$app/forms';

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

	let { data } = $props<{
		user?: { username: string };
		servers: ServerData[];
	}>();

	let loggingOut = $state(false);
	let activeFilter = $state<'all' | 'vanilla' | 'fabric'>('all');

	let filteredServers = $derived(
		activeFilter === 'all'
			? data.servers
			: data.servers.filter((s: ServerData) => s.type.toLowerCase() === activeFilter)
	);

	function setFilter(filter: 'all' | 'vanilla' | 'fabric') {
		activeFilter = filter;
	}

	function getStatusColor(status?: string) {
		switch (status?.toLowerCase()) {
			case 'running':
				return 'text-green-500';
			case 'stopped':
				return 'text-muted-foreground';
			default:
				return 'text-yellow-500';
		}
	}
</script>

<div class="min-h-screen bg-background">
	<header class="border-b">
		<div class="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
			<div>
				<h1 class="text-xl font-medium text-foreground">Servers</h1>
				<p class="text-sm text-muted-foreground">Manage your Minecraft servers</p>
			</div>
			<div class="flex items-center gap-4">
				{#if data.user}
					<span class="text-sm text-muted-foreground">{data.user.username}</span>
				{/if}
				<form
					method="POST"
					action="?/logout"
					use:enhance={() => {
						loggingOut = true;
						return async ({ update }) => {
							await update();
							loggingOut = false;
						};
					}}
				>
					<Button type="submit" variant="ghost" size="sm" disabled={loggingOut}>
						{loggingOut ? 'Logging out...' : 'Log out'}
					</Button>
				</form>
			</div>
		</div>
	</header>

	<main class="mx-auto max-w-4xl px-6 py-8">
		{#if data.servers.length === 0}
			<div class="flex min-h-[60vh] flex-col items-center justify-center text-center">
				<div class="mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-muted">
					<ServerIcon class="h-6 w-6 text-muted-foreground" />
				</div>
				<h2 class="text-lg font-medium text-foreground">No servers yet</h2>
				<p class="mt-1 text-sm text-muted-foreground">
					Create your first Minecraft server to get started
				</p>
				<Button class="mt-6" href="/dashboard/servers/new">
					<Plus class="mr-2 h-4 w-4" />
					Create Server
				</Button>
			</div>
		{:else}
			<div class="mb-6 flex items-center justify-between">
				<p class="text-sm text-muted-foreground">
					{filteredServers.length} server{filteredServers.length === 1 ? '' : 's'}
				</p>
				<Button size="sm" href="/dashboard/servers/new">
					<Plus class="mr-2 h-4 w-4" />
					New Server
				</Button>
			</div>

			<!-- Filter Tags -->
			<div class="mb-8 flex items-center gap-2">
				<button
					type="button"
					onclick={() => setFilter('all')}
					class="rounded-full border px-4 py-1.5 text-sm transition-colors {activeFilter === 'all'
						? 'border-primary bg-primary text-primary-foreground'
						: 'border-input bg-background text-foreground hover:bg-muted'}"
				>
					All
				</button>
				<button
					type="button"
					onclick={() => setFilter('vanilla')}
					class="rounded-full border px-4 py-1.5 text-sm transition-colors {activeFilter ===
					'vanilla'
						? 'border-primary bg-primary text-primary-foreground'
						: 'border-input bg-background text-foreground hover:bg-muted'}"
				>
					Vanilla
				</button>
				<button
					type="button"
					onclick={() => setFilter('fabric')}
					class="rounded-full border px-4 py-1.5 text-sm transition-colors {activeFilter ===
					'fabric'
						? 'border-primary bg-primary text-primary-foreground'
						: 'border-input bg-background text-foreground hover:bg-muted'}"
				>
					Fabric
				</button>
			</div>

			<!-- Server List -->
			<div class="space-y-6">
				{#each filteredServers as server}
					<a
						href="/dashboard/servers/{server.id}"
						class="group block border-b border-border pb-6 transition-colors last:border-0"
					>
						<div class="space-y-3">
							<h3 class="text-xl font-medium text-foreground">{server.name}</h3>
							<p class="text-sm leading-relaxed text-muted-foreground">
								Minecraft server running on port {server.port} with {server.min_mem}MB - {server.max_mem}MB
								memory allocation
							</p>
							<div class="flex items-center gap-2 pt-1">
								<Badge variant="secondary" class="text-xs">
									{server.type}
								</Badge>
								<Badge variant="outline" class="text-xs">v{server.version}</Badge>
							</div>
							<div class="flex items-center gap-4 pt-2 text-xs text-muted-foreground">
								<div class="flex items-center gap-1.5">
									<Circle class="h-3 w-3 {getStatusColor(server.status)}" fill="currentColor" />
									<span class="capitalize">{server.status || 'Unknown'}</span>
								</div>
								<span>Port {server.port}</span>
							</div>
						</div>
					</a>
				{/each}
			</div>
		{/if}
	</main>
</div>
