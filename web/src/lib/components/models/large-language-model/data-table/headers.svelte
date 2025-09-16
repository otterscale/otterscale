<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import { type LargeLangeageModel } from '../protobuf.svelte';

	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		gpu_cache,
		kv_cache,
		requests,
		uptime,
		relation,
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
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet gpu_cache(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.gpu_cache()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet kv_cache(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.gpu_cache()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet requests(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.requests()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet uptime(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.uptime()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet relation()}{/snippet}
