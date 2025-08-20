<script lang="ts">
	import { page } from '$app/state';
	import { Workload } from '$lib/components/applications/workload';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';

	breadcrumb.set({ parents: [], current: dynamicPaths.settings(page.params.scope) });
</script>

{#if $currentKubernetes}
	{@const scopeUuid = $currentKubernetes.scopeUuid}
	{@const facilityName = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}
	{@const applicationName = page.params.application_name ?? ''}

	<Workload {scopeUuid} {facilityName} {namespace} {applicationName} />
{/if}
