<script lang="ts">
	import { Metals } from '$lib/components/machines';
	import { machinesMetalPath, machinesPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parents: [machinesPath], current: machinesMetalPath });
</script>

{#if $activeScope}
	<Metals />
{/if}
