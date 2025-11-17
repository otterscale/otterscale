<script lang="ts" module>
	import type { VirtualMachine_Restore } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		virtualMachineRestore,
		scope,
		reloadManager
	}: {
		virtualMachineRestore: VirtualMachine_Restore;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {virtualMachineRestore} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
