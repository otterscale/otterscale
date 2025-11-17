<script lang="ts" module>
	import type { User, User_Key } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
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
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {key} {user} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
