<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Construction } from '$lib/components/construction';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.account(),
			url: resolve('/(auth)/scope/[scope]/account', { scope: page.params.scope! }),
		},
	]);
</script>

<Construction />
