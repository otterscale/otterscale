<script lang="ts" module>
	import type { SubvolumeGroup } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		poolName: poolName,
		usage: usage,
		mode: mode,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<SubvolumeGroup>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<SubvolumeGroup>)}
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet poolName(row: Row<SubvolumeGroup>)}
	<Badge variant="outline">
		{row.original.poolName}
	</Badge>
{/snippet}

{#snippet mode(row: Row<SubvolumeGroup>)}
	<Badge variant="outline">
		{row.original.mode}
	</Badge>
{/snippet}

{#snippet usage(row: Row<SubvolumeGroup>)}
	<Progress.Root
		numerator={Number(row.original.usedBytes)}
		denominator={Number(row.original.usedBytes)}
	>
		{#snippet detail({ numerator, denominator })}
			{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(
				numerator / (1024 * 1024)
			)}
			{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(
				denominator / (1024 * 1024)
			)}
			<span>
				{numeratorValue}
				{numeratorUnit}
			</span>
			<span>/</span>
			<span>
				{denominatorValue}
				{denominatorUnit}
			</span>
		{/snippet}
		{#snippet ratio({ numerator, denominator })}
			{((numerator * 100) / denominator).toFixed(2)}%
		{/snippet}
	</Progress.Root>
{/snippet}

{#snippet createTime(row: Row<SubvolumeGroup>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}
