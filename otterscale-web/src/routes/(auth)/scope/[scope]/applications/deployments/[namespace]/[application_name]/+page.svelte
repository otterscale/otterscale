<script lang="ts" module>
	import { page } from '$app/state';
	import { Workload } from '$lib/components/applications/workload';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	breadcrumb.set({
		parents: [
			dynamicPaths.applications(page.params.scope),
			dynamicPaths.applicationsDeployments(page.params.scope)
		],
		current: { title: `${page.params.namespace} / ${page.params.application_name}`, url: '' }
	});
</script>

{#if $currentKubernetes}
	{@const scopeUuid = $currentKubernetes.scopeUuid}
	{@const facilityName = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}
	{@const applicationName = page.params.application_name ?? ''}

	<Workload {scopeUuid} {facilityName} {namespace} {applicationName} />
{/if}
