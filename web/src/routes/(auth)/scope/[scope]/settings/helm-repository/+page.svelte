<script lang="ts">
	import { page } from '$app/state';
	import { HelmRepository } from '$lib/components/settings/helm-repository';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.settings(page.params.scope)],
		current: { title: m.helm_repository(), url: '' },
	});
</script>

<HelmRepository />
