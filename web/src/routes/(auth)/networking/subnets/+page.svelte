<script lang="ts">
	import { resolve } from '$app/paths';
	import { Subnet } from '$lib/components/network/subnet';
	import { m } from '$lib/paraglide/messages';
	import { activeScope, breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.networking(),
			url: resolve('/(auth)/networking'),
		},
		{
			title: m.subnets(),
			url: resolve('/(auth)/networking/subnets'),
		},
	]);
</script>

{#if $activeScope}
	<Subnet />
{/if}
