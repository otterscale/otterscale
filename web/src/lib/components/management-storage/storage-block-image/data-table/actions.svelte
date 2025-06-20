<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionCopy from './action-copy.svelte';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import { ActionSnapshot } from './action-snapshot';
	import type { BlockImage } from './types';

	let { blockImage }: { blockImage: BlockImage } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Views</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionSnapshot { blockImage } />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {blockImage} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionCopy {blockImage} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {blockImage} />
	</Layout.ActionItem>
</Layout.Actions>
