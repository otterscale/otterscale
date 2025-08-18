<script lang="ts" module>
	import type { User } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from './action-delete.svelte';
	import Edit from './action-edit.svelte';
</script>

<script lang="ts">
	let {
		user
	}: {
		user: User;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {user} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {user} />
	</Layout.ActionItem>
</Layout.Actions>
