<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import Delete from '../actions/delete.svelte';
	import Edit from '../actions/edit.svelte';
</script>

<script lang="ts">
	let {
		image
	}: {
		image: Image;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {image} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {image} />
	</Layout.ActionItem>
</Layout.Actions>
