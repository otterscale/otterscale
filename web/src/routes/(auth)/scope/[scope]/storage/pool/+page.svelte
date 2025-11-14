<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Pool } from '$lib/components/storage/pool';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';
</script>

<script lang="ts">
	breadcrumbs.set([
		{
			title: m.storage(),
			url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! })
		},
		{
			title: m.pool(),
			url: resolve('/(auth)/scope/[scope]/storage/pool', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope}
	<Pool scope={page.params.scope!} />
{/key}
