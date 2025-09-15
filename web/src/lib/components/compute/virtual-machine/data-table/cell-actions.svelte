<script lang="ts" module>
	import Clone from './action-clone.svelte';
	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
	import PauseResume from './action-pause-resume.svelte';
	import Snapshot from './action-snapshot.svelte';
	import StartStop from './action-start-stop.svelte';
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
		<PauseResume {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<StartStop {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<StartupScript {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Snapshot {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Edit {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Clone {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {virtualMachine} />
	</Layout.ActionItem>
</Layout.Actions>
