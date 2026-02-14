<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Construction } from '$lib/components/construction';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.databases(),
			url: resolve('/(auth)/scope/[scope]/databases', { scope: page.params.scope! })
		},
		{
			title: m.relational(),
			url: resolve('/(auth)/scope/[scope]/databases/relational', { scope: page.params.scope! })
		}
	]);
</script>

<Construction />
