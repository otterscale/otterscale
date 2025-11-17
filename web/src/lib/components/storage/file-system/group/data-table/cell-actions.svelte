<script lang="ts" module>
	import type { SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';
	import { m } from '$lib/paraglide/messages';

	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		subvolumeGroup,
		scope,
		volume,
		reloadManager
	}: {
		subvolumeGroup: SubvolumeGroup;
		scope: string;
		volume: string;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>{m.actions()}</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {subvolumeGroup} {scope} {volume} {reloadManager} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {subvolumeGroup} {scope} {volume} {reloadManager} />
	</Layout.ActionItem>
</Layout.Actions>
