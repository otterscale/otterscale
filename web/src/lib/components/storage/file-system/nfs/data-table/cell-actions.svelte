<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
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
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Grant {subvolume} {scope} {volume} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Revoke {subvolume} {scope} {volume} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Edit {subvolume} {scope} {volume} {group} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {subvolume} {scope} {volume} {group} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
