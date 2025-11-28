<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import PowerOff from './action-power-off.svelte';
	import Remove from './action-remove.svelte';
</script>

<script lang="ts">
	let {
		machine,
		reloadManager
	}: {
		machine: Machine;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem
		disabled={(machine.status.toLowerCase() !== 'ready' &&
			machine.status.toLowerCase() !== 'releasing') ||
			!!machine.workloadAnnotations['juju-is-controller']}
	>
		<Remove {machine} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem
		disabled={machine.powerState.toLowerCase() !== 'on' ||
			machine.status.toLowerCase() === 'commissioning' ||
			machine.status.toLowerCase() === 'testing' ||
			machine.status.toLowerCase() === 'deploying' ||
			!!machine.workloadAnnotations['juju-is-controller']}
	>
		<PowerOff {machine} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
