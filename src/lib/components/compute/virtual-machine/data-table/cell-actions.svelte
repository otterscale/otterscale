<script lang="ts" module>
	import type { VirtualMachine } from '$lib/api/instance/v1/instance_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Clone from './action-clone.svelte';
	import Delete from './action-delete.svelte';
	import Migrate from './action-migrate.svelte';
	import PauseResume from './action-pause-resume.svelte';
	import Restart from './action-restart.svelte';
	import Restore from './action-restore.svelte';
	import Snapshot from './action-snapshot.svelte';
	import StartStop from './action-start-stop.svelte';
</script>

<script lang="ts">
	let {
		virtualMachine,
		scope,
		reloadManager
	}: { virtualMachine: VirtualMachine; scope: string; reloadManager: ReloadManager } = $props();

	let open = $state(false);
	function close() {
		open = false;
	}
</script>

<Actions.List bind:open>
	<Actions.Label>{m.actions()}</Actions.Label>
	<Actions.Separator />
	<Actions.Item>
		<PauseResume {virtualMachine} {scope} />
	</Actions.Item>
	<Actions.Item>
		<StartStop {virtualMachine} {scope} />
	</Actions.Item>
	<Actions.Item>
		<Restart {virtualMachine} {scope} />
	</Actions.Item>
	<Actions.Separator />
	<Actions.Item>
		<Snapshot {virtualMachine} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Restore {virtualMachine} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Separator />
	<Actions.Item>
		<Clone {virtualMachine} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Migrate {virtualMachine} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {virtualMachine} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
