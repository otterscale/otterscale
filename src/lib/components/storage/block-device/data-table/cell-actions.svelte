<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		image,
		scope,
		reloadManager
	}: {
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
	<Actions.Item>
		<Edit {image} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
	<Actions.Item>
		<Delete {image} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
