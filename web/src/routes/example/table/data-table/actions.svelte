<script lang="ts" module>
	import type { TableRow } from './type';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Writable } from 'svelte/store';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		row,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		row: TableRow;
		data: Writable<TableRow[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>Edit</Layout.ActionItem>
	<Layout.ActionItem>Delete</Layout.ActionItem>
</Layout.Actions>
