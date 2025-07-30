<script lang="ts">
	import { storageClusterPath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [storagePath], current: storageClusterPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
