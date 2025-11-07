<script lang="ts">
	import { resolve } from '$app/paths';
	import { PackageRepository } from '$lib/components/settings/package-repository';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/configuration'),
		},
		{
			title: m.package_repository(),
			url: resolve('/(auth)/configuration/package-repository'),
		},
	]);
</script>

<PackageRepository />
