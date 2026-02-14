<script lang="ts">
	import { resolve } from '$app/paths';
	import { SingleSignOn } from '$lib/components/settings/single-sign-on';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.configuration(),
			url: resolve('/(auth)/configuration')
		},
		{
			title: m.single_sign_on(),
			url: resolve('/(auth)/configuration/single-sign-on')
		}
	]);
</script>

<SingleSignOn />
