<script lang="ts">
	import { Metal } from '$lib/components/machines/metal/data-table/metal';
	import { machinesMetalPath, machinesPath } from '$lib/path';
	import { breadcrumb, scopeLoading } from '$lib/stores';

	breadcrumb.set({ parent: machinesPath, current: machinesMetalPath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	<Metal />
{/if}
