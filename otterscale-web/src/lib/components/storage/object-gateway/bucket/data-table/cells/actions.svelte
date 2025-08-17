<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from '../actions/delete.svelte';
	import Edit from '../actions/edit.svelte';
</script>

<script lang="ts">
	let {
		bucket
	}: {
		bucket: Bucket;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {bucket} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {bucket} />
	</Layout.ActionItem>
</Layout.Actions>
