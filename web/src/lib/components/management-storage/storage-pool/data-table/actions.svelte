<script lang="ts">
	import type { Pool } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';

	let { pool }: { pool: Pool } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {pool} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {pool} />
	</Layout.ActionItem>
</Layout.Actions>
