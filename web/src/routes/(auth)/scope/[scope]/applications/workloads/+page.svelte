<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Workloads } from '$lib/components/applications/workloads';
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
	]);
</script>

{#if $currentKubernetes}
	{@const scope = $currentKubernetes.scope}
	{@const facility = $currentKubernetes.name}

	<Workloads {scope} {facility} />
{/if}
