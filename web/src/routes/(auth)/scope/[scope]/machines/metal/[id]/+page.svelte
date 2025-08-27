<script lang="ts">
	import { page } from '$app/state';
	import { Metal } from '$lib/components/machines/metal/index';
	import { dynamicPaths } from '$lib/path';
	import { activeScope, breadcrumb } from '$lib/stores';

	breadcrumb.set({
		parents: [
			dynamicPaths.machines(page.params.scope),
			dynamicPaths.machinesMetal(page.params.scope)
		],
		current: { title: page.params.id!, url: '' }
	});
</script>

{#if $activeScope}
	<Metal />
{/if}
