<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import {
		Card,
		CardContent,
		CardDescription,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Select, SelectContent, SelectItem, SelectTrigger } from '$lib/components/ui/select';
	import { enhance } from '$app/forms';
	import { Boxes, Info, ArrowLeft, Loader2 } from '@lucide/svelte';
	import type { ServerType } from './+page.server';

	const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3001';

	let { data, form } = $props<{
		data: { user: { username: string }; serverTypes: ServerType[] };
		form?: { error?: string; values?: Record<string, string | null> };
	}>();

	let creating = $state(false);
	let selectedTypeId = $state(form?.values?.serverType || '');
	let selectedVersion = $state(form?.values?.version || '');
	let nameLength = $state((form?.values?.name || '').length);
	let availableVersions = $state<string[]>([]);
	let loadingVersions = $state(false);
	let versionsError = $state('');

	// Get the selected server type details
	let selectedType = $derived(data.serverTypes.find((t: ServerType) => t.id === selectedTypeId));

	async function fetchVersions(typeId: string) {
		if (!typeId) {
			availableVersions = [];
			return;
		}

		loadingVersions = true;
		versionsError = '';
		selectedVersion = '';

		try {
			const response = await fetch(`${API_URL}/api/v1/versions/${typeId}`);
			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ error: 'Failed to fetch versions' }));
				throw new Error(errorData.error || 'Failed to fetch versions');
			}
			const result = await response.json();
			availableVersions = result.versions;
		} catch (error) {
			versionsError = error instanceof Error ? error.message : 'Failed to load versions';
			availableVersions = [];
		} finally {
			loadingVersions = false;
		}
	}

	function handleTypeChange() {
		fetchVersions(selectedTypeId);
	}

	function handleNameInput(event: Event) {
		const target = event.target as HTMLInputElement;
		nameLength = target.value.length;
	}
</script>

<div class="min-h-screen bg-background">
	<header class="border-b">
		<div class="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
			<div class="flex items-center gap-4">
				<Button variant="ghost" size="sm" href="/dashboard">
					<ArrowLeft class="mr-2 h-4 w-4" />
					Back
				</Button>
				<div>
					<h1 class="text-xl font-medium text-foreground">Create Server</h1>
					<p class="text-sm text-muted-foreground">Set up a new Minecraft server</p>
				</div>
			</div>
			<div class="flex items-center gap-4">
				{#if data.user}
					<span class="text-sm text-muted-foreground">{data.user.username}</span>
				{/if}
			</div>
		</div>
	</header>

	<main class="mx-auto max-w-6xl px-6 py-8">
		<div class="flex min-h-[60vh] items-center justify-center">
			<Card class="w-full max-w-md">
				<CardHeader class="text-center">
					<div
						class="mx-auto mb-4 flex h-12 w-12 items-center justify-center rounded-full bg-muted"
					>
						<Boxes class="h-6 w-6 text-muted-foreground" />
					</div>
					<CardTitle class="text-lg font-medium">Create Server</CardTitle>
					<CardDescription>Configure your new Minecraft server</CardDescription>
				</CardHeader>
				<CardContent>
					{#if form?.error}
						<div
							class="mb-4 rounded-lg border border-destructive/50 bg-destructive/10 px-4 py-3 text-sm text-destructive"
						>
							{form.error}
						</div>
					{/if}

					<form
						method="POST"
						action="?/create"
						use:enhance={() => {
							creating = true;
							return async ({ update }) => {
								await update();
								creating = false;
							};
						}}
						class="space-y-6"
					>
						<!-- Server Name -->
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
								maxlength={100}
								value={form?.values?.name || ''}
								oninput={handleNameInput}
								class="h-11"
								required
							/>
							<div class="flex justify-end">
								<span class="text-xs text-muted-foreground">{nameLength} / 100</span>
							</div>
						</div>

						<!-- Server Type -->
						<div class="space-y-2">
							<Label for="serverType">
								Server Type
								<span class="text-destructive">*</span>
							</Label>
							<input type="hidden" name="serverType" value={selectedTypeId} />
							<Select type="single" bind:value={selectedTypeId} onValueChange={handleTypeChange} required>
								<SelectTrigger class="w-full">
									{#if selectedType}
										{selectedType.name}
									{:else}
										<span class="text-muted-foreground">Choose a server type</span>
									{/if}
								</SelectTrigger>
								<SelectContent>
									{#each data.serverTypes as type (type.id)}
										<SelectItem value={type.id} label={type.name} />
									{/each}
								</SelectContent>
							</Select>
							{#if selectedType}
								<p class="text-xs text-muted-foreground">{selectedType.description}</p>
							{/if}
						</div>

						<!-- Server Version -->
						<div class="space-y-2">
							<Label for="version">
								Server Version
								<span class="text-destructive">*</span>
							</Label>
							<input type="hidden" name="version" bind:value={selectedVersion} />
							<Select type="single" bind:value={selectedVersion} disabled={!selectedType || loadingVersions} required>
								<SelectTrigger class="w-full">
									{#if loadingVersions}
										<span class="flex items-center gap-2 text-muted-foreground">
											<Loader2 class="h-4 w-4 animate-spin" />
											Loading versions...
										</span>
									{:else if selectedVersion}
										{selectedVersion}
									{:else}
										<span class="text-muted-foreground">
											{selectedType ? 'Choose a version' : 'Select a server type first'}
										</span>
									{/if}
								</SelectTrigger>
								<SelectContent>
									{#each availableVersions as version (version)}
										<SelectItem value={version} label={version} />
									{/each}
								</SelectContent>
							</Select>
							{#if versionsError}
								<p class="text-xs text-destructive">{versionsError}</p>
							{/if}
						</div>

						<!-- Divider -->
						<div class="border-t"></div>

						<!-- Port -->
						<div class="space-y-2">
							<Label for="port">Port</Label>
							<Input
								id="port"
								name="port"
								type="number"
								placeholder="25565"
								min={1024}
								max={65535}
								value={form?.values?.port || ''}
								class="h-11"
							/>
							<p class="text-xs text-muted-foreground">Auto-assigned if left empty</p>
						</div>

						<!-- Memory Allocation -->
						<div class="space-y-3">
							<Label>Memory Allocation (GB)</Label>
							<div class="grid grid-cols-2 gap-4">
								<div class="space-y-2">
									<Label for="minMemory" class="text-xs text-muted-foreground">Minimum</Label>
									<Input
										id="minMemory"
										name="minMemory"
										type="number"
										placeholder="2"
										min={1}
										value={form?.values?.minMemory || '2'}
										class="h-11"
									/>
								</div>
								<div class="space-y-2">
									<Label for="maxMemory" class="text-xs text-muted-foreground">Maximum</Label>
									<Input
										id="maxMemory"
										name="maxMemory"
										type="number"
										placeholder="4"
										min={1}
										max={128}
										value={form?.values?.maxMemory || '4'}
										class="h-11"
									/>
								</div>
							</div>
						</div>

						<!-- Info Note -->
						<div class="flex items-start gap-2 rounded-lg bg-muted px-3 py-2">
							<Info class="mt-0.5 h-4 w-4 flex-shrink-0 text-muted-foreground" />
							<p class="text-xs text-muted-foreground">These settings can be changed later</p>
						</div>

						<!-- Actions -->
						<div class="flex items-center justify-between gap-4 pt-2">
							<Button type="button" variant="ghost" href="/dashboard" disabled={creating}>
								Cancel
							</Button>
							<Button type="submit" disabled={creating} class="min-w-[140px]">
								{creating ? 'Creating...' : 'Create Server'}
							</Button>
						</div>
					</form>
				</CardContent>
			</Card>
		</div>
	</main>
</div>
