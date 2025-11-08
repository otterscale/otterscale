<script lang="ts">
	import { resolve } from '$app/paths';
	import { Subscription } from '$lib/components/settings/subscription';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs, premiumTier } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/configuration')
		},
		{
			title: m.subscription(),
			url: resolve('/(auth)/configuration/subscription')
		}
	]);
</script>

{#key $premiumTier.level}
	<Subscription {premiumTier} />
{/key}
