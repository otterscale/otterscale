<script lang="ts">
	import { page } from '$app/state';
	import { Store } from '$lib/components/applications/store';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.applications(page.params.scope)],
		current: dynamicPaths.applicationsStore(page.params.scope),
	});
</script>

{#if $activeScope}
	<Store />
{/if}
