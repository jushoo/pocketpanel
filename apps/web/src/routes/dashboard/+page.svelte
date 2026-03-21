<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { enhance } from "$app/forms";
	
	let { data } = $props();
	let loggingOut = $state(false);
</script>

<div class="min-h-screen flex items-center justify-center bg-background">
	<div class="text-center">
		<h1 class="text-xl font-medium text-foreground">Dashboard</h1>
		{#if data.user}
			<p class="text-sm text-muted-foreground mt-2">Welcome, {data.user.username}</p>
		{:else}
			<p class="text-sm text-muted-foreground mt-2">Welcome to PocketPanel</p>
		{/if}
		
		<form method="POST" action="?/logout" use:enhance={() => {
			loggingOut = true;
			return async ({ update }) => {
				await update();
				loggingOut = false;
			};
		}} class="mt-8">
			<Button type="submit" variant="outline" disabled={loggingOut}>
				{loggingOut ? "Logging out..." : "Log out"}
			</Button>
		</form>
	</div>
</div>
