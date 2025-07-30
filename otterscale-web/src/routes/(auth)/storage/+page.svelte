<script lang="ts">
	import { homePath, storagePath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [homePath], current: storagePath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
