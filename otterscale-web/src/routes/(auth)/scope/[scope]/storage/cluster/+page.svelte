<script lang="ts" module>
	import { page } from '$app/state';
	import { Cluster } from '$lib/components/storage/cluster';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storageCluster(page.params.scope)
	});
</script>

{#if $activeScope}
	<Cluster />
{/if}
