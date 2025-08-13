<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		selectedFacility,
		row,
		images: data = $bindable()
	}: {
		selectedScopeUuid: string;
		selectedFacility: string;
		row: Row<Image>;
		images: Writable<Image[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionItem>
		<Edit {selectedScopeUuid} {selectedFacility} image={row.original} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} {selectedFacility} image={row.original} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
