<script lang="ts" module>
	import type { User, User_Key } from '$lib/api/storage/v1/storage_pb';
	import { Actions } from '$lib/components/custom/data-table/core';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		key,
		user,
		scope,
		reloadManager
	}: {
		key: User_Key;
		user: User;
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
		<Delete {key} {user} {scope} {reloadManager} closeActions={close} />
	</Actions.Item>
</Actions.List>
