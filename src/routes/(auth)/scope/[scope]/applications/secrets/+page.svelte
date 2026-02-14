<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Secrets } from '$lib/components/applications/secrets';
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
			title: m.secrets(),
			url: resolve('/(auth)/scope/[scope]/applications/secrets', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<Secrets scope={page.params.scope!} />
{/key}
