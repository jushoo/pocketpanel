<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { Input } from "$lib/components/ui/input";
	import { Label } from "$lib/components/ui/label";

	let username = $state("");
	let password = $state("");
	let loading = $state(false);
	let error = $state("");

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = "";
		loading = true;

		await new Promise((resolve) => setTimeout(resolve, 1000));
		
		if (username && password) {
			console.log("Login attempted", { username });
		} else {
			error = "Please enter both username and password";
		}
		
		loading = false;
	}
</script>

<div class="min-h-screen w-full flex items-center justify-center bg-background">
	<div class="w-full max-w-sm px-6">
		<div class="mb-12 text-center">
			<div class="inline-flex items-center justify-center w-12 h-12 rounded-xl bg-muted border border-border mb-6">
				<svg 
					viewBox="0 0 24 24" 
					class="w-6 h-6 text-foreground"
					fill="none" 
					stroke="currentColor" 
					stroke-width="1.5"
				>
					<rect x="3" y="3" width="18" height="18" rx="2" />
					<path d="M9 12h6M12 9v6" />
				</svg>
			</div>
			<h1 class="text-xl font-medium text-foreground tracking-tight">PocketPanel</h1>
			<p class="mt-2 text-sm text-muted-foreground">Admin Login</p>
		</div>

		<form onsubmit={handleSubmit} class="space-y-5">
			<div class="space-y-2">
				<Label for="username" class="text-sm font-normal text-muted-foreground">Username</Label>
				<Input
					id="username"
					type="text"
					placeholder="admin"
					bind:value={username}
					class="h-11"
				/>
			</div>

			<div class="space-y-2">
				<Label for="password" class="text-sm font-normal text-muted-foreground">Password</Label>
				<Input
					id="password"
					type="password"
					placeholder="••••••••"
					bind:value={password}
					class="h-11"
				/>
			</div>

			{#if error}
				<p class="text-sm text-destructive">{error}</p>
			{/if}

			<Button
				type="submit"
				disabled={loading}
				class="w-full h-11"
			>
				{loading ? "Signing in..." : "Sign in"}
			</Button>
		</form>

		<p class="mt-8 text-center text-xs text-muted-foreground">
			Secure admin access only
		</p>
	</div>
</div>
