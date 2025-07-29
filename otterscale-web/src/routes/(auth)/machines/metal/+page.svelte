<script lang="ts">
	import { machinesMetalPath, machinesPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: machinesPath, current: machinesMetalPath });
</script>

{#if $activeScope}
	current scope: {$activeScope.uuid}
{/if}
