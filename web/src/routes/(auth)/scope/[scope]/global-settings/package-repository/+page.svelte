<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { PackageRepository } from '$lib/components/settings/package-repository';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.package_repository(),
			url: resolve('/(auth)/scope/[scope]/global-settings/package-repository', { scope: page.params.scope! }),
		},
	]);
</script>

<PackageRepository />
