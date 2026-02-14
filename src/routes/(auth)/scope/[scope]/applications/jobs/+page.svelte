<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Jobs } from '$lib/components/applications/jobs';
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
			title: m.jobs(),
			url: resolve('/(auth)/scope/[scope]/applications/jobs', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<Jobs scope={page.params.scope!} />
{/key}
