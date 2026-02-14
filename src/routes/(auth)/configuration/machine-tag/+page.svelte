<script lang="ts">
	import { resolve } from '$app/paths';
	import { MachineTag } from '$lib/components/settings/machine-tag';
	import { m } from '$lib/paraglide/messages';
	import { breadcrumbs } from '$lib/stores';

	// Set breadcrumbs navigation
	breadcrumbs.set([
		{
			title: m.configuration(),
			url: resolve('/(auth)/configuration')
		},
		{
			title: m.machine_tag(),
			url: resolve('/(auth)/configuration/machine-tag')
		}
	]);
</script>

<MachineTag />
