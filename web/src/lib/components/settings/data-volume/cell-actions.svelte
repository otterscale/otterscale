<script lang="ts" module>
	import Delete from './action-delete.svelte';
	import Extend from './action-extend.svelte';
	// import Migrate from './action-migrate.svelte';
	// import PauseResume from './action-pause-resume.svelte';
	// import Restart from './action-restart.svelte';
	// import Restore from './action-restore.svelte';
	// import Snapshot from './action-snapshot.svelte';
	// import StartStop from './action-start-stop.svelte';

	import type { DataVolume } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
</script>

<script lang="ts">
	let { dataVolume }: { dataVolume: DataVolume } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Extend {dataVolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {dataVolume} />
	</Layout.ActionItem>
	<!-- <Layout.ActionItem>
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
 -->
</Layout.Actions>
