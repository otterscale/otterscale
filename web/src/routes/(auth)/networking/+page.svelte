<script lang="ts">
	import { resolve } from '$app/paths';
	import Dashboard from '$lib/components/network/dashboard/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.networking(),
			url: resolve('/(auth)/networking')
		}
	]);
</script>

<Dashboard />
