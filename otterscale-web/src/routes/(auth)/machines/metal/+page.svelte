<script lang="ts">
	import { Metals } from '$lib/components/machines';
	import { machinesMetalPath, machinesPath } from '$lib/path';
	import { breadcrumb, scopeLoading } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({ parent: machinesPath, current: machinesMetalPath });
</script>

{#if $scopeLoading}
	loading scope
{:else}
	<Metals />
{/if}
