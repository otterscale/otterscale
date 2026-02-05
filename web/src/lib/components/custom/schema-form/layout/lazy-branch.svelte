<script lang="ts">
	import { onMount, tick } from 'svelte';

	// Accept the component to render and all props to pass through
	let { component: Component, ...restProps }: { component: any; [key: string]: any } = $props();

	let isVisible = $state(false);

	onMount(async () => {
		// Wait for the next microtask to clear the call stack and ensure pending state updates are applied
		await tick();
		isVisible = true;
	});
</script>

{#if isVisible}
	<Component {...restProps} />
{:else}
	<!-- Placeholder while waiting for async render -->
{/if}
