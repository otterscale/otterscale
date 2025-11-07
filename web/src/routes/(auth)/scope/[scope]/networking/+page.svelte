<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Dashboard from '$lib/components/network/dashboard/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.networking(),
			url: resolve('/(auth)/scope/[scope]/networking', { scope: page.params.scope! }),
		},
	]);
</script>

<Dashboard />
