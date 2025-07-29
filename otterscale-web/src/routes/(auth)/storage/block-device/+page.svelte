<script lang="ts">
	import { storageBlockDevicePath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb, loadingScopes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: storagePath, current: storageBlockDevicePath });
</script>

{#if $loadingScopes}
	loading scope
{:else}
	current scope: {$activeScope.uuid}
{/if}
