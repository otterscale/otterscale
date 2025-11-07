<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { MachineTag } from '$lib/components/settings/machine-tag';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.settings(),
			url: resolve('/(auth)/scope/[scope]/global-settings', { scope: page.params.scope! }),
		},
		{
			title: m.machine_tag(),
			url: resolve('/(auth)/scope/[scope]/global-settings/machine-tag', { scope: page.params.scope! }),
		},
	]);
</script>

<MachineTag />
