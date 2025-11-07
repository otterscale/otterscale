<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { HelmRepository } from '$lib/components/settings/helm-repository';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.helm_repository(),
			url: resolve('/(auth)/scope/[scope]/global-settings/helm-repository', { scope: page.params.scope! }),
		},
	]);
</script>

<HelmRepository />
