<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Subscription } from '$lib/components/settings/subscription';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, premiumTier } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.subscription(),
			url: resolve('/(auth)/scope/[scope]/global-settings/subscription', { scope: page.params.scope! }),
		},
	]);
</script>

{#key $premiumTier.level}
	<Subscription {premiumTier} />
{/key}
