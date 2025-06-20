<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionCopy from './action-copy.svelte';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { Snapshot } from './types';

	let { blockImagesnapshot }: { blockImagesnapshot: Snapshot } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {blockImagesnapshot} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionCopy {blockImagesnapshot} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {blockImagesnapshot} />
	</Layout.ActionItem>
</Layout.Actions>
