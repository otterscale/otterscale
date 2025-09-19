<script lang="ts" module>
	import Delete from './action-delete.svelte';

	import type { VirtualMachineSnapshot } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { snapshot }: { snapshot: VirtualMachineSnapshot } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete virtualMachineSnapshot={snapshot} />
	</Layout.ActionItem>
</Layout.Actions>
