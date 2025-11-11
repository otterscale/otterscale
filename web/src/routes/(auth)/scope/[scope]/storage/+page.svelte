<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Dashboard from '$lib/components/storage/dashboard/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.storage(),
			url: resolve('/(auth)/scope/[scope]/storage', { scope: page.params.scope! })
		}
	]);
</script>

<Dashboard scope={page.params.scope!} />
