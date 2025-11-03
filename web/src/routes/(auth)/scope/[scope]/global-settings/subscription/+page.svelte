<script lang="ts">
	import { page } from '$app/state';
	import { Subscription } from '$lib/components/settings/subscription';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb, premiumTier } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.globalSettings(page.params.scope)],
		current: { title: m.subscription(), url: '' },
	});
</script>

{#key $premiumTier.level}
	<Subscription {premiumTier} />
{/key}
