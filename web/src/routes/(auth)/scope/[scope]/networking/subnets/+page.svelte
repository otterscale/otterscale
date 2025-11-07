<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Subnet } from '$lib/components/network/subnet';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.networking(),
			url: resolve('/(auth)/scope/[scope]/networking', { scope: page.params.scope! }),
		},
		{
			title: m.subnets(),
			url: resolve('/(auth)/scope/[scope]/networking/subnets', { scope: page.params.scope! }),
		},
	]);
</script>

{#if $activeScope}
	<Subnet />
{/if}
