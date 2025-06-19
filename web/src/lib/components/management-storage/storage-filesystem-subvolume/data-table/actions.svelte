<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import ActionEdit from './action-edit.svelte';
	import type { Subvolume } from './types';

	let { subvolume }: { subvolume: Subvolume } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionEdit {subvolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {subvolume} />
	</Layout.ActionItem>
</Layout.Actions>
