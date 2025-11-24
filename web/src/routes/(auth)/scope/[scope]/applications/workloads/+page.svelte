<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Workloads } from '$lib/components/applications/workloads';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	breadcrumbs.set([
		{
			title: m.applications(),
			url: resolve('/(auth)/scope/[scope]/applications', { scope: page.params.scope! })
		},
		{
			title: m.workloads(),
			url: resolve('/(auth)/scope/[scope]/applications/workloads', { scope: page.params.scope! })
		}
	]);
	console.log(page.data);
</script>

{#key page.params.scope!}
	<Workloads scope={page.params.scope!} />
{/key}
