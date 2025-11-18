<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import Dashboard from '$lib/components/machines/dashboard/index.svelte';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.machines(),
			url: resolve('/(auth)/machines')
		}
	]);
</script>

{#key page.params.scope!}
	<Dashboard scope={page.params.scope!} />
{/key}
