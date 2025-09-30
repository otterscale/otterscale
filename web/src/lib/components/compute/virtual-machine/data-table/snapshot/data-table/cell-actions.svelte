<script lang="ts" module>
	import Delete from './action-delete.svelte';

	import type { VirtualMachine_Snapshot } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachineSnapshot }: { virtualMachineSnapshot: VirtualMachine_Snapshot } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {virtualMachineSnapshot} />
	</Layout.ActionItem>
</Layout.Actions>
