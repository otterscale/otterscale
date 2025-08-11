<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker,
		name,
		poolName,
		usage
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
	<Progress.Root {numerator} {denominator}>
		{#snippet ratio({ numerator, denominator })}
			{Math.round((numerator * 100) / denominator)}%
		{/snippet}
	</Progress.Root>
{/snippet}
