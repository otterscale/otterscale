<script lang="ts">
	import { Metal } from '$lib/components/machines/metal/data-table/metal';
	import { machinesMetalPath, machinesPath } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	breadcrumb.set({ parents: [machinesPath], current: machinesMetalPath });
</script>

{#if $activeScope}
	<Metal />
{/if}
