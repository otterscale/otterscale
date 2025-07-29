<script lang="ts">
	import { storageClusterPath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb, loadingScopes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: storagePath, current: storageClusterPath });
</script>

{#if $loadingScopes}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}
