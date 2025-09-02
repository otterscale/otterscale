<script lang="ts" module>
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { formatBigNumber } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import { type LargeLangeageModel } from '../protobuf.svelte';

	export const cells = {
		row_picker,
		name,
		version,
		parameters,
		accuracy,
		speed,
		architecture,
		requests,
		uptime,
	};
</script>

{#snippet row_picker(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet version(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-start">
		{row.original.version}
	</Layout.Cell>
{/snippet}

{#snippet parameters(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.parameters}
	</Layout.Cell>
{/snippet}

{#snippet accuracy(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.metrics.accuracy}
	</Layout.Cell>
{/snippet}

{#snippet speed(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.metrics.speed}
	</Layout.Cell>
{/snippet}

{#snippet architecture(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-start">
		{row.original.architecture}
	</Layout.Cell>
{/snippet}

{#snippet requests(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.usageStats.requests)}
	</Layout.Cell>
{/snippet}

{#snippet uptime(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.usageStats.uptime)}
	</Layout.Cell>
{/snippet}
