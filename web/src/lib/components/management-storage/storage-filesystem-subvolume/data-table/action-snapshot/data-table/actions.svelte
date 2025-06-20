<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { Snapshot } from './types';

	let { snapshot }: { snapshot: Snapshot } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {snapshot} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {snapshot} />
	</Layout.ActionItem>
</Layout.Actions>
