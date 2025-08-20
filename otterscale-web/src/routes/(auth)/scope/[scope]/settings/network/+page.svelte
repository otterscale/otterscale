<script lang="ts">
	import { page } from '$app/state';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';
	import { Network } from '$lib/components/settings/network';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: dynamicPaths.settingsNetwork(page.params.scope)
	});
</script>

{#if $activeScope}
	<Network />
{/if}
