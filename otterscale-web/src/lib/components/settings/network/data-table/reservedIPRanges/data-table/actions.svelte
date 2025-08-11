<script lang="ts" module>
	import type { Network, Network_IPRange } from '$lib/api/network/v1/network_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import Delete from './delete.svelte';
	import Update from './update.svelte';
</script>

<script lang="ts">
	let {
		row,
		networks = $bindable()
	}: {
		row: Row<Network_IPRange>;
		networks: Writable<Network[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Update ipRange={row.original} bind:networks />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Delete ipRange={row.original} bind:networks />
	</Layout.ActionItem>
</Layout.Actions>
