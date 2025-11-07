<script lang="ts">
	import { resolve } from '$app/paths';
	import { HelmRepository } from '$lib/components/settings/helm-repository';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/configuration'),
		},
		{
			title: m.helm_repository(),
			url: resolve('/(auth)/configuration/helm-repository'),
		},
	]);
</script>

<HelmRepository />
