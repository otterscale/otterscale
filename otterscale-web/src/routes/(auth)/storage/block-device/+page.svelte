<script lang="ts">
	import { storageBlockDevicePath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [storagePath], current: storageBlockDevicePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
