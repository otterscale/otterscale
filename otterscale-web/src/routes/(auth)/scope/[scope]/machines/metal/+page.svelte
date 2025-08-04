<script lang="ts">
	import { page } from '$app/state';
	import { Metals } from '$lib/components/machines';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.machines(page.params.scope)],
		current: dynamicPaths.machinesMetal(page.params.scope)
	});
</script>

{#if $activeScope}
	<Metals />
{/if}
