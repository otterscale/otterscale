<script lang="ts" module>
	import type { Image } from '$lib/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacityV2 as formatCapacity } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cells/actions.svelte';
	import { Snapshots } from './cells/snapshots';

	export const cells = {
		_row_picker,
		name,
		poolName,
		usage,
		snapshots,
		actions
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
			{Progress.formatRatio(numerator, denominator)}
		{/snippet}
		{#snippet detail({ numerator, denominator })}
			{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(numerator)}
			{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(denominator)}
			{numeratorValue}
			{numeratorUnit}/{denominatorValue}
			{denominatorUnit}
		{/snippet}
	</Progress.Root>
{/snippet}

{#snippet snapshots(row: Row<Image>)}
	<Snapshots image={row.original} />
{/snippet}

{#snippet actions(row: Row<Image>)}
	<Actions image={row.original} />
{/snippet}
