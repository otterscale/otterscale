<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Update from './action-update.svelte';
</script>

<script lang="ts">
	let {
		ipRange,
		reloadManager
	}: {
		ipRange: Network_IPRange;
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
		<Update {ipRange} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {ipRange} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
