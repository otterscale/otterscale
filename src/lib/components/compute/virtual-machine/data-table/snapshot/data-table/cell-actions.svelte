<script lang="ts" module>
	import type { VirtualMachine_Snapshot } from '$lib/api/instance/v1/instance_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		virtualMachineSnapshot,
		scope,
		reloadManager
	}: {
		virtualMachineSnapshot: VirtualMachine_Snapshot;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();

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
	<Actions.Item>
		<Delete {virtualMachineSnapshot} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
