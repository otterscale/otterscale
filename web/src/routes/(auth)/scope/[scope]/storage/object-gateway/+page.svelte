<script lang="ts" module>
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { ObjectDateway } from '$lib/components/storage/object-gateway';
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
			title: m.object_gateway(),
			url: resolve('/(auth)/scope/[scope]/storage/object-gateway', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<ObjectDateway scope={page.params.scope!} />
{/key}
