<script lang="ts" module>
	import Clone from './action-clone.svelte';
	import Delete from './action-delete.svelte';
	import Migrate from './action-migrate.svelte';
	import PauseResume from './action-pause-resume.svelte';
	import Restart from './action-restart.svelte';
	import Restore from './action-restore.svelte';
	import Snapshot from './action-snapshot.svelte';
	import StartStop from './action-start-stop.svelte';

	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { virtualMachine }: { virtualMachine: VirtualMachine } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<PauseResume {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<StartStop {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Restart {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Snapshot {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Restore {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Clone {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Migrate {virtualMachine} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {virtualMachine} />
	</Layout.ActionItem>
</Layout.Actions>
