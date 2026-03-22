<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import { ArrowLeft, Circle, Play, Square, Loader2, Terminal } from '@lucide/svelte';
	import { enhance } from '$app/forms';
	import { formatMemory } from '$lib/utils.js';

	let { data, form } = $props<{
		server: {
			id: number;
			name: string;
			type: string;
			version: string;
			min_mem: number;
			max_mem: number;
			port: number;
		};
		running: boolean;
		pid?: number;
		form?: {
			error?: string;
			success?: boolean;
			message?: string;
		};
	}>();

	let starting = $state(false);
	let stopping = $state(false);
	let showConsole = $state(false);
	let consoleLines = $state<string[]>([]);

	async function fetchConsoleHistory() {
		if (!data.running) return;
		
		try {
			const response = await fetch(`/api/proxy/servers/${data.server.id}/console?lines=50`);
			if (response.ok) {
				const data = await response.json();
				consoleLines = data.lines || [];
			}
		} catch {
			// Ignore errors
		}
	}

	function toggleConsole() {
		showConsole = !showConsole;
		if (showConsole && data.running) {
			fetchConsoleHistory();
		}
	}

	function getStatusColor(running: boolean) {
		return running ? 'text-green-500' : 'text-muted-foreground';
	}

	function getStatusText(running: boolean) {
		return running ? 'Running' : 'Stopped';
	}
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
		<div class="mb-8 flex items-start justify-between">
			<div>
				<h1 class="text-xl font-medium text-foreground">{data.server.name}</h1>
				<p class="mt-1 text-sm text-muted-foreground">Server configuration and details</p>
			</div>
			
			<!-- Start/Stop Buttons -->
			<div class="flex items-center gap-2">
				{#if data.running}
					<form
						method="POST"
						action="?/stop"
						use:enhance={() => {
							stopping = true;
							return async ({ update }) => {
								await update();
								stopping = false;
							};
						}}
					>
						<Button variant="outline" size="sm" disabled={stopping}>
							{#if stopping}
								<Loader2 class="mr-2 h-4 w-4 animate-spin" />
								Stopping...
							{:else}
								<Square class="mr-2 h-4 w-4" />
								Stop Server
							{/if}
						</Button>
					</form>
				{:else}
					<form
						method="POST"
						action="?/start"
						use:enhance={() => {
							starting = true;
							return async ({ update }) => {
								await update();
								starting = false;
							};
						}}
					>
						<Button size="sm" disabled={starting}>
							{#if starting}
								<Loader2 class="mr-2 h-4 w-4 animate-spin" />
								Starting...
							{:else}
								<Play class="mr-2 h-4 w-4" />
								Start Server
							{/if}
						</Button>
					</form>
				{/if}
				
				<Button variant="ghost" size="sm" onclick={toggleConsole}>
					<Terminal class="mr-2 h-4 w-4" />
					Console
				</Button>
			</div>
		</div>

		<!-- Error/Success Messages -->
		{#if form?.error}
			<div class="mb-6 rounded-md border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive">
				{form.error}
			</div>
		{/if}
		{#if form?.success}
			<div class="mb-6 rounded-md border border-green-500/50 bg-green-500/10 px-4 py-3 text-sm text-green-600">
				{form.message}
			</div>
		{/if}

		<!-- Console Output -->
		{#if showConsole}
			<div class="mb-8">
				<div class="mb-2 flex items-center justify-between">
					<span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Console Output</span>
					<Button variant="ghost" size="sm" onclick={fetchConsoleHistory}>
						Refresh
					</Button>
				</div>
				<div class="rounded-md border border-border bg-black/5 p-4 font-mono text-xs">
					{#if data.running}
						{#if consoleLines.length === 0}
							<p class="text-muted-foreground">Waiting for output...</p>
						{:else}
							{#each consoleLines as line}
								<div class="text-foreground">{line}</div>
							{/each}
						{/if}
					{:else}
						<p class="text-muted-foreground">Server is not running. Start the server to see console output.</p>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Server Info - List Style -->
		<div class="space-y-6">
			<!-- Name -->
			<div class="pb-6 border-b border-border">
				<div class="mb-2 flex items-center justify-between">
					<span class="text-xs font-medium uppercase tracking-wider text-muted-foreground">Name</span>
				</div>
				<p class="text-lg font-medium text-foreground">{data.server.name}</p>
			</div>

			<!-- Type & Version Row -->
			<div class="pb-6 border-b border-border">
				<span class="mb-2 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Version</span>
				<div class="flex items-center gap-3">
					<p class="text-foreground">v{data.server.version}</p>
					<Badge variant="secondary" class="text-xs">{data.server.type}</Badge>
				</div>
			</div>

			<!-- Resources -->
			<div class="pb-6 border-b border-border">
				<span class="mb-3 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Resources</span>
				<div class="grid grid-cols-2 gap-6">
					<div>
						<p class="text-xs text-muted-foreground">Memory</p>
						<p class="mt-1 font-medium text-foreground">
							{formatMemory(data.server.min_mem)} - {formatMemory(data.server.max_mem)}
						</p>
					</div>
					<div>
						<p class="text-xs text-muted-foreground">Port</p>
						<p class="mt-1 font-medium text-foreground">{data.server.port}</p>
					</div>
				</div>
			</div>

			<!-- Status -->
			<div class="pb-6 border-b border-border">
				<span class="mb-2 block text-xs font-medium uppercase tracking-wider text-muted-foreground">Status</span>
				<div class="flex items-center gap-2">
					<Circle class="h-3 w-3 {getStatusColor(data.running)}" fill="currentColor" />
					<span class="text-sm {data.running ? 'text-green-600' : 'text-muted-foreground'}">
						{getStatusText(data.running)}
					</span>
					{#if data.pid}
						<span class="text-xs text-muted-foreground">(PID: {data.pid})</span>
					{/if}
				</div>
			</div>
		</div>
	</main>
</div>
