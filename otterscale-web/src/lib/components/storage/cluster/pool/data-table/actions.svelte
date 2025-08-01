<script lang="ts" module>
	import type { Pool } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Edit from './edit.svelte';
</script>

<script lang="ts">
	let {
		selectedScopeUuid,
		row,
		pools: data = $bindable()
	}: {
		selectedScopeUuid: string;
		row: Row<Pool>;
		pools: Writable<Pool[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Edit {selectedScopeUuid} pool={row.original} bind:data />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete {selectedScopeUuid} pool={row.original} bind:data />
	</Layout.ActionItem>
</Layout.Actions>
