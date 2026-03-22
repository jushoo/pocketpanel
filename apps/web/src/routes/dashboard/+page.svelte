<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Badge } from '$lib/components/ui/badge';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Plus, Server } from '@lucide/svelte';
	import { enhance } from '$app/forms';

	let { data } = $props();
	let loggingOut = $state(false);
</script>

<div class="min-h-screen bg-background">
	<header class="border-b">
		<div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
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

	<main class="mx-auto max-w-6xl px-6 py-8">
		{#if data.servers.length === 0}
			<div class="flex min-h-[60vh] items-center justify-center">
				<Card class="w-full max-w-md">
					<CardHeader class="text-center">
						<div
							class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-muted"
						>
							<Server class="h-6 w-6 text-muted-foreground" />
						</div>
						<CardTitle class="text-lg font-medium">No servers yet</CardTitle>
						<CardDescription>Create your first Minecraft server to get started</CardDescription>
					</CardHeader>
					<CardContent>
						<Button class="w-full" size="lg" href="/dashboard/servers/new">
							<Plus class="mr-2 h-4 w-4" />
							Create Server
						</Button>
					</CardContent>
				</Card>
			</div>
		{:else}
			<div class="mb-6 flex items-center justify-between">
				<p class="text-sm text-muted-foreground">{data.servers.length} server{data.servers.length === 1 ? '' : 's'}</p>
				<Button size="sm" href="/dashboard/servers/new">
					<Plus class="mr-2 h-4 w-4" />
					New Server
				</Button>
			</div>
			<div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
				{#each data.servers as server}
					<Card>
						<CardHeader class="pb-2">
							<div class="flex items-center justify-between">
								<CardTitle class="text-base font-medium">{server.name}</CardTitle>
								<Badge variant="secondary" class="text-xs">{server.type}</Badge>
							</div>
							<CardDescription class="text-xs">Port {server.port}</CardDescription>
						</CardHeader>
						<CardContent>
							<div class="flex items-center justify-between text-sm">
								<span class="text-muted-foreground">v{server.version}</span>
								<span class="text-muted-foreground">{server.min_mem}MB - {server.max_mem}MB</span>
							</div>
						</CardContent>
					</Card>
				{/each}
			</div>
		{/if}
	</main>
</div>
