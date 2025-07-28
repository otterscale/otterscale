<script lang="ts" module>
	import type { Pool } from '$gen/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScope,
		selectedFacility,
		row,
		pools: data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		row: Row<Pool>;
		pools: Writable<Pool[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {selectedScope} {selectedFacility} pool={row.original} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScope} {selectedFacility} pool={row.original} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
