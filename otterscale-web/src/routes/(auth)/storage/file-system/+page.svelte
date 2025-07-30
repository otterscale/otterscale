<script lang="ts">
	import { storageFileSystemPath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [storagePath], current: storageFileSystemPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
