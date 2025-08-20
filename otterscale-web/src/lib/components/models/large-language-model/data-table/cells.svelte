<script lang="ts" module>
	import { RowPickers } from '$lib/components/custom/data-table';
	import type { Row } from '@tanstack/table-core';
	import { type LargeLangeageModel } from '../protobuf.svelte';
	import { formatBigNumber } from '$lib/formatter';

	export const cells = {
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

{#snippet row_picker(row: Row<LargeLangeageModel>)}
	<RowPickers.Cell {row} />
{/snippet}

{#snippet name(row: Row<LargeLangeageModel>)}
	{row.original.name}
{/snippet}

{#snippet version(row: Row<LargeLangeageModel>)}
	{row.original.version}
{/snippet}

{#snippet parameters(row: Row<LargeLangeageModel>)}
	<div class="flex justify-end">
		{row.original.parameters}
	</div>
{/snippet}

{#snippet accuracy(row: Row<LargeLangeageModel>)}
	<div class="flex justify-end">
		{row.original.metrics.accuracy}
	</div>
{/snippet}

{#snippet speed(row: Row<LargeLangeageModel>)}
	<span class="flex justify-end">
		{row.original.metrics.speed}
	</span>
{/snippet}

{#snippet architecture(row: Row<LargeLangeageModel>)}
	{row.original.architecture}
{/snippet}

{#snippet requests(row: Row<LargeLangeageModel>)}
	<span class="flex justify-end">
		{formatBigNumber(row.original.usageStats.requests)}
	</span>
{/snippet}

{#snippet uptime(row: Row<LargeLangeageModel>)}
	<span class="flex justify-end">
		{formatBigNumber(row.original.usageStats.requests)}
	</span>
{/snippet}
