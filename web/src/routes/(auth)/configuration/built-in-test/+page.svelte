<script lang="ts">
	import { resolve } from '$app/paths';
	import { BuiltInSelfTest } from '$lib/components/settings/built-in-self-test/index';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/configuration'),
		},
		{
			title: m.built_in_test(),
			url: resolve('/(auth)/configuration/built-in-test'),
		},
	]);
</script>

<BuiltInSelfTest />
