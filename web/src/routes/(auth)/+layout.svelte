<script lang="ts">
	import { onMount, type Snippet } from 'svelte';

	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { isAuthenticated } from '$lib/stores';

	let { children }: { children: Snippet } = $props();

	onMount(() => {
		const unsubscribe = isAuthenticated.subscribe((value) => {
			if (!value) {
				goto(resolve('/'));
			}
		});

		return () => unsubscribe();
	});
</script>

{@render children()}
