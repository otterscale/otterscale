<script lang="ts" module>
	import { page } from '$app/state';
	import { OSD } from '$lib/components/storage/osd';
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
		<OSD />
	{/key}
{/if}
