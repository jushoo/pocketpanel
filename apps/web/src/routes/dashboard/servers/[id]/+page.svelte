<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { ArrowLeft, Circle } from '@lucide/svelte';
	import { formatMemory } from '$lib/utils.js';

	let { data } = $props<{
		server: {
			id: number;
			name: string;
			type: string;
			version: string;
			min_mem: number;
			max_mem: number;
			port: number;
		};
	}>();

	let server = $derived(data.server);
</script>

<div class="min-h-screen bg-background">
	<header class="border-b">
		<div class="mx-auto flex max-w-4xl items-center justify-between px-6 py-4">
			<div class="flex items-center gap-4">
				<Button variant="ghost" size="sm" href="/dashboard">
					<ArrowLeft class="mr-2 h-4 w-4" />
					Back to servers
				</Button>
			</div>
		</div>
	</header>

	<main class="mx-auto max-w-4xl px-6 py-8">
		<!-- Title -->
		<div class="mb-8">
			<h1 class="text-xl font-medium text-foreground">{server.name}</h1>
			<p class="mt-1 text-sm text-muted-foreground">Server configuration and details</p>
		</div>

		<!-- Server Info - List Style -->
		<div class="space-y-6">
			<!-- Name -->
			<div class="pb-6 border-b border-border">
				<div class="mb-2 flex items-center justify-between">
					<span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Name</span>
				</div>
				<p class="text-lg font-medium text-foreground">{server.name}</p>
			</div>

			<!-- Type & Version Row -->
			<div class="pb-6 border-b border-border">
				<span class="mb-2 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Version</span>
				<div class="flex items-center gap-3">
					<p class="text-foreground">v{server.version}</p>
					<Badge variant="secondary" class="text-xs">{server.type}</Badge>
				</div>
			</div>

			<!-- Resources -->
			<div class="pb-6 border-b border-border">
				<span class="mb-3 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Resources</span>
				<div class="grid grid-cols-2 gap-6">
					<div>
						<p class="text-xs text-muted-foreground">Memory</p>
						<p class="mt-1 font-medium text-foreground">
							{formatMemory(server.min_mem)} - {formatMemory(server.max_mem)}
						</p>
					</div>
					<div>
						<p class="text-xs text-muted-foreground">Port</p>
						<p class="mt-1 font-medium text-foreground">{server.port}</p>
					</div>
				</div>
			</div>

			<!-- Status -->
			<div class="pb-6 border-b border-border">
				<span class="mb-2 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Status</span>
				<div class="flex items-center gap-2">
					<Circle class="h-3 w-3 text-yellow-500" fill="currentColor" />
					<span class="text-sm text-muted-foreground">Unknown</span>
				</div>
			</div>
		</div>
	</main>
</div>
