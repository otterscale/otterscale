<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Construction } from '$lib/components/construction';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.registry(),
			url: resolve('/(auth)/scope/[scope]/registry', { scope: page.params.scope! })
		}
	]);
</script>

<Construction />
