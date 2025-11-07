<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Workload } from '$lib/components/applications/workload';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, currentKubernetes } from '$lib/stores';
</script>

<script lang="ts">
	breadcrumbs.set([
		{
			title: m.applications(),
			url: resolve('/(auth)/scope/[scope]/applications', { scope: page.params.scope! }),
		},
		{
			title: m.workloads(),
			url: resolve('/(auth)/scope/[scope]/applications/workloads', { scope: page.params.scope! }),
		},
		{
			title: page.params.namespace!,
			url: resolve('/(auth)/scope/[scope]/applications/workloads/[namespace]', {
				scope: page.params.scope!,
				namespace: page.params.namespace!,
			}),
		},
		{
			title: page.params.application_name!,
			url: resolve('/(auth)/scope/[scope]/applications/workloads/[namespace]/[application_name]', {
				scope: page.params.scope!,
				namespace: page.params.namespace!,
				application_name: page.params.application_name!,
			}),
		},
	]);
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}
	{@const namespace = page.params.namespace ?? ''}
	{@const applicationName = page.params.application_name ?? ''}

	<Workload {scope} {facility} {namespace} {applicationName} />
{/if}
