<script lang="ts">
	import { page } from '$app/state';
	import { MachineTag } from '$lib/components/settings/machine-tag';
	import { m } from '$lib/paraglide/messages';
	import { dynamicPaths } from '$lib/path';
	import { breadcrumb } from '$lib/stores';

	// Set breadcrumb navigation
	breadcrumb.set({
		parents: [dynamicPaths.globalSettings(page.params.scope)],
		current: { title: m.machine_tag(), url: '' },
	});
</script>

<MachineTag />
