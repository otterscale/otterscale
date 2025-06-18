<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { Pool } from './types';

	let { pool }: { pool: Pool } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {pool} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {pool} />
	</Layout.ActionItem>
</Layout.Actions>
