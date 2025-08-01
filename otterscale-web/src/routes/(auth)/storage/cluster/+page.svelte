<script lang="ts" module>
	import { Cluster } from '$lib/components/storage/cluster';
	import { storageClusterPath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [storagePath], current: storageClusterPath });
</script>

{#if $activeScope}
	<Cluster />
{/if}
