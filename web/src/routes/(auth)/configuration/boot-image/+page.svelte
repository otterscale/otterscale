<script lang="ts">
	import { resolve } from '$app/paths';
	import { BootImage } from '$lib/components/settings/boot-image';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/configuration')
		},
		{
			title: m.boot_image(),
			url: resolve('/(auth)/configuration/boot-image')
		}
	]);
</script>

<BootImage />
