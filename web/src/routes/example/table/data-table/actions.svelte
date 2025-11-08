<script lang="ts" module>
	import type { Writable } from 'svelte/store';

	import type { TableRow } from './type';

	import * as Layout from '$lib/components/custom/data-table/layout';
</script>

<script lang="ts">
	let {
		data = $bindable()
	}: {
		data: Writable<TableRow[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>Edit</Layout.ActionItem>
	<Layout.ActionItem>Delete</Layout.ActionItem>
</Layout.Actions>
