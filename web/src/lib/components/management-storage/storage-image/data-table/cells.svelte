<script lang="ts" module>
	import type { Image } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		poolName: poolName,
		usage: usage
	};
</script>

{#snippet _row_picker(row: Row<Image>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Image>)}
	{row.original.name}
{/snippet}

{#snippet poolName(row: Row<Image>)}
	<Badge variant="outline">{row.original.poolName}</Badge>
{/snippet}

{#snippet usage(row: Row<Image>)}
	{@const denominator = Number(row.original.quotaBytes)}
	{@const numerator = Number(row.original.usedBytes)}
	<div class="flex justify-end">
		<Progress.Root {numerator} {denominator}>
			{#snippet ratio({ numerator, denominator })}
				{Math.round((numerator * 100) / denominator)}%
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{numerator}/{denominator}
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}
