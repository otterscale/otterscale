<script lang="ts">
	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from './action-delete.svelte';

	let {
		snapshot
	}: {
		snapshot: Subvolume_Snapshot;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {snapshot} />
	</Layout.ActionItem>
</Layout.Actions>
