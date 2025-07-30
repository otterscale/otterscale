<script lang="ts">
	import { storageObjectGatewayPath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [storagePath], current: storageObjectGatewayPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
