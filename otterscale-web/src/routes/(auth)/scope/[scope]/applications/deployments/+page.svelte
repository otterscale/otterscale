<script lang="ts" module>
	import { page } from '$app/state';
	import { Workloads } from '$lib/components/applications/workloads';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	breadcrumb.set({
		parents: [dynamicPaths.applications(page.params.scope)],
		current: dynamicPaths.applicationsDeployments(page.params.scope)
	});
</script>

{#if $currentKubernetes}
	{@const scopeUuid = $currentKubernetes.scopeUuid}
	{@const facilityName = $currentKubernetes.name}

	<Workloads {scopeUuid} {facilityName} />
{/if}
