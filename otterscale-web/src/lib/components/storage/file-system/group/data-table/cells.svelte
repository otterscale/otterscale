<script lang="ts" module>
	import type { SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		poolName,
		usage,
		mode,
		createTime,
		actions
	};
</script>

{#snippet row_picker(row: Row<SubvolumeGroup>)}
	<RowPicker {row} />
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
	<div class="flex justify-end">
		<Progress.Root
			numerator={Number(row.original.usedBytes)}
			denominator={Number(row.original.usedBytes)}
			highIsGood={false}
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{@const { value: numeratorValue, unit: numeratorUnit } = formatCapacity(numerator)}
				{@const { value: denominatorValue, unit: denominatorUnit } = formatCapacity(denominator)}
				{numeratorValue}
				{numeratorUnit}
				/
				{denominatorValue}
				{denominatorUnit}
			{/snippet}
		</Progress.Root>
	</div>
{/snippet}

{#snippet createTime(row: Row<SubvolumeGroup>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}

{#snippet actions(row: Row<SubvolumeGroup>)}
	<Actions subvolumeGroup={row.original} />
{/snippet}
