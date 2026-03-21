<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { ModeWatcher } from 'mode-watcher';
	import ModeToggle from '$lib/components/mode-toggle.svelte';

	let { children } = $props();
</script>

<ModeWatcher />

<div class="fixed top-4 right-4 z-50">
	<ModeToggle />
</div>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>PocketPanel - Admin</title>
</svelte:head>

{@render children()}
