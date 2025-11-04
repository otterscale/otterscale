<script lang="ts">
	import { page } from '$app/state';
	import { BootImage } from '$lib/components/settings/boot-image';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.globalSettings(page.params.scope)],
		current: { title: m.boot_image(), url: '' },
	});
</script>

<BootImage />
