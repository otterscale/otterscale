<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Extensions } from '$lib/components/settings/extensions';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! })
		},
		{
			title: m.extensions(),
			url: resolve('/(auth)/scope/[scope]/settings/extensions', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<Extensions scope={page.params.scope!} />
{/key}
