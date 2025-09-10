<script lang="ts" module>
	import { page } from '$app/state';
	import { DataVolume } from '$lib/components/compute/datavolume';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	breadcrumb.set({
		parents: [dynamicPaths.compute(page.params.scope), dynamicPaths.computeVirtualMachine(page.params.scope)],
		current: { title: `${page.params.namespace} / ${page.params.virtual_machine_name}`, url: '' },
	});
</script>

{#if $currentKubernetes}
	{@const scopeUuid = $currentKubernetes.scopeUuid}
	{@const facilityName = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}
	{@const virtualMachineName = page.params.virtual_machine_name ?? ''}
	<DataVolume {scopeUuid} {facilityName} {namespace} {virtualMachineName} />
{/if}
