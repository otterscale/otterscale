<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [], current: dynamicPaths.settings(page.params.scope) });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
