<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import type { Row } from '@tanstack/table-core';
	import type { Writable } from 'svelte/store';
	import Add from './add.svelte';
	import PowerOff from './power-off.svelte';
	import Remove from './remove.svelte';
</script>

<script lang="ts">
	let {
		row,
		machines = $bindable()
	}: {
		row: Row<Machine>;
		machines: Writable<Machine[]>;
	} = $props();
</script>

<Layout.Actions>
	<Layout.ActionLabel>Actions</Layout.ActionLabel>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<Add machine={row.original} bind:machines />
	</Layout.ActionItem>
	<Layout.ActionItem>
		<Remove machine={row.original} bind:machines />
	</Layout.ActionItem>
	<Layout.ActionSeparator />
	<Layout.ActionItem>
		<PowerOff machine={row.original} bind:machines />
	</Layout.ActionItem>
</Layout.Actions>
