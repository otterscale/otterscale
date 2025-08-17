<script lang="ts" module>
	import type { SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from '../actions/delete.svelte';
	import Edit from '../actions/edit.svelte';
</script>

<script lang="ts">
	let {
		subvolumeGroup
	}: {
		subvolumeGroup: SubvolumeGroup;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {subvolumeGroup} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {subvolumeGroup} />
	</Layout.ActionItem>
</Layout.Actions>
