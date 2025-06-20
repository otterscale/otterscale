<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import ActionDelete from './action-delete.svelte';
	import { ActionSnapshot } from './action-snapshot';
	import ActionEdit from './action-edit.svelte';
	import type { Subvolume } from './types';

	let { subvolume }: { subvolume: Subvolume } = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<ActionSnapshot {subvolume}/>
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionEdit {subvolume} />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<ActionDelete {subvolume} />
	</Layout.ActionItem>
</Layout.Actions>
