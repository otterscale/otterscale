<script lang="ts" module>
	import type { Pool } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		pool
	}: {
		pool: Pool;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Edit {pool} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {pool} />
	</Layout.ActionItem>
</Layout.Actions>
