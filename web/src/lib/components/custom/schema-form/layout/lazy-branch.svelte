<script lang="ts">
	import { onMount } from 'svelte';

	// Accept the component to render and all props to pass through
	let { component: Component, ...restProps }: { component: any; [key: string]: any } = $props();

	let isVisible = $state(false);

	onMount(() => {
		// Key: setTimeout 0 defers execution to the next Macrotask
		// This lets the browser think "the current task is done", clearing the Call Stack
		setTimeout(() => {
			isVisible = true;
		}, 0);
	});
</script>

{#if isVisible}
	<Component {...restProps} />
{:else}
	<!-- Placeholder while waiting for async render -->
{/if}
