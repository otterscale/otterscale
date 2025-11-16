<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { BuiltInSelfTest } from '$lib/components/settings/built-in-self-test/index';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/settings', { scope: page.params.scope! })
		},
		{
			title: m.built_in_test(),
			url: resolve('/(auth)/scope/[scope]/settings/built-in-test', { scope: page.params.scope! })
		}
	]);
</script>

{#key page.params.scope!}
	<BuiltInSelfTest scope={page.params.scope!} />
{/key}
