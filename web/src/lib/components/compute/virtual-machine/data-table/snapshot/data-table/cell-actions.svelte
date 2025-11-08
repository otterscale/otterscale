<script lang="ts" module>
	import type { VirtualMachine_Snapshot } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
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
