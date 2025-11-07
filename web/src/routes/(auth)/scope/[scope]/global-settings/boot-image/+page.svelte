<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { BootImage } from '$lib/components/settings/boot-image';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.boot_image(),
			url: resolve('/(auth)/scope/[scope]/global-settings/boot-image', { scope: page.params.scope! }),
		},
	]);
</script>

<BootImage />
