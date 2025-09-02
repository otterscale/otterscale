<script lang="ts">
	import { page } from '$app/state';
	import Dashboard from '$lib/components/storage/dashboard/index.svelte';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.storage(page.params.scope) });
</script>

{#if $activeScope}
	<Dashboard scope={$activeScope} />
{/if}
