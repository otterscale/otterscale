<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
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

	const scope = $derived(machine.workloadAnnotations['juju-machine-id']?.split('-machine-')[0]);

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Actions.List bind:open>
	<Actions.Label>
		{m.actions()}
	</Actions.Label>
	<Actions.Separator />
	<Actions.Item disabled={!scope}>
		<Remove {machine} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item
		disabled={machine.powerState.toLowerCase() !== 'on' ||
			machine.status.toLowerCase() === 'commissioning' ||
			machine.status.toLowerCase() === 'testing' ||
			machine.status.toLowerCase() === 'deploying' ||
			!!machine.workloadAnnotations['juju-is-controller']}
	>
		<PowerOff {machine} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
