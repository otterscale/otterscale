<script lang="ts" module>
	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from './action-delete.svelte';
</script>

<script lang="ts">
	let {
		key
	}: {
		key: User_Key;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Delete {key} />
	</Layout.ActionItem>
</Layout.Actions>
