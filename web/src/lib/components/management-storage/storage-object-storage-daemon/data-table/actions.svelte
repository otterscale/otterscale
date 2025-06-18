<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import { ActionViewDevice } from './action-view-device';
	import type { ObjectStorageDaemon } from './types';

	let { osd }: { osd: ObjectStorageDaemon } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>View</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionViewDevice />
	</Layout.ActionItem>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {osd} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {osd} />
	</Layout.ActionItem>
</Layout.Actions>
