<script lang="ts">
	import { page } from '$app/state';
	import { VirtualMachine } from '$lib/components/compute/virtual-machine';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.compute(page.params.scope)],
		current: dynamicPaths.computeVirtualMachine(page.params.scope),
	});
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}

	<VirtualMachine {scope} {facility} {namespace} />
{/if}
