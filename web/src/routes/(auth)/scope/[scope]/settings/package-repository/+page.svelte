<script lang="ts">
	import { page } from '$app/state';
	import { PackageRepository } from '$lib/components/settings/package-repository';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.package_repository(), url: '' },
	});
</script>

<PackageRepository />
