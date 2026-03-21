<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { enhance } from '$app/forms';

	let { form } = $props();
	let loading = $state(false);
</script>

<div class="flex min-h-screen w-full items-center justify-center bg-background">
	<div class="w-full max-w-sm px-6">
		<div class="mb-12 text-center">
			<div
				class="mb-6 inline-flex h-12 w-12 items-center justify-center rounded-xl border border-border bg-muted"
			>
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
			<h1 class="text-xl font-medium tracking-tight text-foreground">PocketPanel</h1>
			<p class="mt-2 text-sm text-muted-foreground">Admin Login</p>
		</div>

		<form
			method="POST"
			action="?/login"
			use:enhance={() => {
				loading = true;
				return async ({ update }) => {
					await update();
					loading = false;
				};
			}}
			class="space-y-5"
		>
			<div class="space-y-2">
				<Label for="username" class="text-sm font-normal text-muted-foreground">Username</Label>
				<Input
					id="username"
					name="username"
					type="text"
					placeholder="admin"
					class="h-11"
					required
				/>
			</div>

			<div class="space-y-2">
				<Label for="password" class="text-sm font-normal text-muted-foreground">Password</Label>
				<Input
					id="password"
					name="password"
					type="password"
					placeholder="••••••••"
					class="h-11"
					required
				/>
			</div>

			{#if form?.error}
				<p class="text-sm text-destructive">{form.error}</p>
			{/if}

			<Button type="submit" disabled={loading} class="h-11 w-full">
				{loading ? 'Signing in...' : 'Sign in'}
			</Button>
		</form>

		<p class="mt-8 text-center text-xs text-muted-foreground">Secure admin access only</p>
	</div>
</div>
