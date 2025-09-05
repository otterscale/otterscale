<script lang="ts">
	import { page } from '$app/state';
	import { Subnet } from '$lib/components/network/subnet';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.networking(page.params.scope)],
		current: dynamicPaths.networkingSubnets(page.params.scope),
	});
</script>

{#if $activeScope}
	<Subnet />
{/if}
