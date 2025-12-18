<script lang="ts" module>
	import type { Image, Image_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Protect from './action-protect.svelte';
	import Rollback from './action-rollback.svelte';
	import Unprotect from './action-unprotect.svelte';
</script>

<script lang="ts">
	let {
		snapshot,
		image,
		scope,
		reloadManager
	}: {
		snapshot: Image_Snapshot;
		image: Image;
		scope: string;
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
		<Rollback {snapshot} {image} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item disabled={snapshot.protected}>
		<Protect {snapshot} {image} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item disabled={!snapshot.protected}>
		<Unprotect {snapshot} {image} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {snapshot} {image} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
