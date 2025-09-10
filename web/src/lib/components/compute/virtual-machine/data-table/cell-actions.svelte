<script lang="ts" module>
	import Snapshot from './action-snapshot.svelte';
	import StartupScript from './action-startup-script.svelte';

	import type { VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Snapshot {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<StartupScript {virtualMachine} />
	</Layout.ActionItem>
</Layout.Actions>
