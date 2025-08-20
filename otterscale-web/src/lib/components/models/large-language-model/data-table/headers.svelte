<script lang="ts">
	import { Layout, RowPickers, Sorter } from '$lib/components/custom/data-table';
	import type { Column, Table } from '@tanstack/table-core';
	import { type LargeLangeageModel } from '../protobuf.svelte';

	export const headers = {
		row_picker,
		name,
		version,
		parameters,
		accuracy,
		speed,
		architecture,
		requests,
		uptime
	};
</script>

{#snippet row_picker(table: Table<LargeLangeageModel>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPickers.Header {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>NAME</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet version(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>VERSION</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet parameters(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>PARAMETERS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet accuracy(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>ACCURACY</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet speed(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SPEED</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet architecture(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>ARCHITECTURE</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet requests(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>REQUESTS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet uptime(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>UPTIME</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
