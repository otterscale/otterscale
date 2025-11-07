<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { SingleSignOn } from '$lib/components/settings/single-sign-on';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.single_sign_on(),
			url: resolve('/(auth)/scope/[scope]/global-settings/single-sign-on', { scope: page.params.scope! }),
		},
	]);
</script>

<SingleSignOn />
