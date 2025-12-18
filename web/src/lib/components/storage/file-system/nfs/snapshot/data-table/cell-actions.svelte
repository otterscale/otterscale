<script lang="ts" module>
	import type { Subvolume, Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		snapshot,
		subvolume,
		scope,
		volume,
		group,
		reloadManager
	}: {
		snapshot: Subvolume_Snapshot;
		subvolume: Subvolume;
		scope: string;
		volume: string;
		group: string;
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
		<Delete {snapshot} {subvolume} {scope} {volume} {group} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
