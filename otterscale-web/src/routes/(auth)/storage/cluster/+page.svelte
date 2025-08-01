<script lang="ts" module>
	import { Cluster } from '$lib/components/storage/cluster';
	import { storageClusterPath, storagePath } from '$lib/path';
	import { breadcrumb, scopeLoading } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: storagePath, current: storageClusterPath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	<Cluster />
{/if}
