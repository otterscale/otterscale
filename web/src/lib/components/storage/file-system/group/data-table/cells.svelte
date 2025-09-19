<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { SubvolumeGroup } from '$lib/api/storage/v1/storage_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		poolName,
		usage,
		mode,
		createTime,
		actions,
	};
</script>

{#snippet row_picker(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-right">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet poolName(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-right">
		<Badge variant="outline">
			{row.original.poolName}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet mode(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-right">
		<Badge variant="outline">
			{row.original.mode}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet usage(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-end">
		{#if row.original.quotaBytes === 0n}
			<span class="text-muted-foreground text-sm">Quota limit is not set</span>
		{:else}
			<Progress.Root
				numerator={Number(row.original.usedBytes)}
				denominator={Number(row.original.quotaBytes)}
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
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-right">
		{#if row.original.createdAt}
			{formatTimeAgo(timestampDate(row.original.createdAt))}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<SubvolumeGroup>)}
	<Layout.Cell class="items-right">
		<Actions subvolumeGroup={row.original} />
	</Layout.Cell>
{/snippet}
