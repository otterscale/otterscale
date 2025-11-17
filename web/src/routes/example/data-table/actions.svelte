<script lang="ts" module>
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { ReloadManager } from '$lib/components/custom/reloader';

	import type { TableRow } from './type';
</script>

<script lang="ts">
	let {
		row,
		reloadManager
	}: {
		row: TableRow;
		reloadManager: ReloadManager;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem
		onclick={() => {
			reloadManager.force();
		}}>Update</Layout.ActionItem
	>
	<Layout.ActionItem>Delete {row.id}</Layout.ActionItem>
</Layout.Actions>
