<script lang="ts" module>
	import type { Subvolume, Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
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
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {snapshot} {subvolume} {scope} {volume} {group} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
