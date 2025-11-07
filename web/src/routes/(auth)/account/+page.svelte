<script lang="ts">
	import { resolve } from '$app/paths';
	import { Construction } from '$lib/components/construction';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.account(),
			url: resolve('/(auth)/account'),
		},
	]);
</script>

<Construction />
