<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './actions-delete.svelte';
	import Edit from './actions-edit.svelte';
	import Grant from './actions-grant-export-access.svelte';
	import Revoke from './actions-revoke-export-access.svelte';
</script>

<script lang="ts">
	let {
		subvolume,
		scope,
		volume,
		group,
		reloadManager
	}: {
		subvolume: Subvolume;
		scope: string;
		volume: string;
		group: string;
		reloadManager: ReloadManager;
	} = $props();

	const exportedClients = $derived(subvolume?.export?.clients ?? []);

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
		<Grant {subvolume} {scope} {volume} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item disabled={exportedClients.length <= 1}>
		<Revoke {subvolume} {scope} {volume} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Edit {subvolume} {scope} {volume} {group} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {subvolume} {scope} {volume} {group} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
