<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import { type LargeLangeageModel } from '../protobuf.svelte';

	import { messages } from './columns';

	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';

	export const headers = {
		row_picker,
		name,
		version,
		parameters,
		accuracy,
		speed,
		architecture,
		requests,
		uptime,
		topology,
	};
</script>

{#snippet row_picker(table: Table<LargeLangeageModel>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages['name']}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
{#snippet version(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages['version']}</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>{messages['parameters']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet accuracy(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['accuracy']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet speed(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['speed']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet architecture(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages['architecture']}</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>{messages['requests']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet uptime(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['uptime']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet topology()}{/snippet}
