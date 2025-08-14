<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { Network } from '$lib/components/settings/network';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.settings(page.params.scope) });
</script>

{#if $activeScope}
	<Network />
{/if}
