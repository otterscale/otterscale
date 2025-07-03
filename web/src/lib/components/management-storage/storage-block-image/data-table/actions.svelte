<script lang="ts">
	import * as Layout from '$lib/components/custom/data-table/layout';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import type { Image } from '$gen/api/storage/v1/storage_pb';

	let {
		selectedScope,
		selectedFacility,
		row,
		data = $bindable()
	}: {
		selectedScope: string;
		selectedFacility: string;
		row: Row<Image>;
		data: Writable<Image[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {selectedScope} {selectedFacility} image={row.original} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScope} {selectedFacility} image={row.original} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
