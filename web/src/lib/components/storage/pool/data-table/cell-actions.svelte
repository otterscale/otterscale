<script lang="ts" module>
	import type { Pool } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		pool,
		scope,
		reloadManager
	}: {
		pool: Pool;
		scope: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Edit {pool} {scope} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {pool} {scope} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
