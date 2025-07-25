<script lang="ts">
	import { storageObjectGatewayPath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb, scopeLoading } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: storagePath, current: storageObjectGatewayPath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}
