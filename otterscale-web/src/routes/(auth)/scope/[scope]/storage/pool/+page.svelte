<script lang="ts" module>
	import { page } from '$app/state';
	import { Pool } from '$lib/components/storage/pool';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.storage(page.params.scope)],
		current: dynamicPaths.storagePool(page.params.scope)
	});
</script>

{#if $activeScope}
	{#key $activeScope.uuid}
		<Pool />
	{/key}
{/if}
