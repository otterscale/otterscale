<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import type { Subvolume } from '$gen/api/storage/v1/storage_pb';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { timestampDate } from '@bufbuild/protobuf/wkt';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		poolName: poolName,
		usage: usage,
		path: path,
		mode: mode,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<Subvolume>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Subvolume>)}
	{row.original.name}
{/snippet}

{#snippet path(row: Row<Subvolume>)}
	<p class="max-w-[200px] overflow-auto text-xs font-light">{row.original.path}</p>
{/snippet}

{#snippet mode(row: Row<Subvolume>)}
	<Badge variant="outline">
		{row.original.mode}
	</Badge>
{/snippet}

{#snippet poolName(row: Row<Subvolume>)}
	<Badge variant="outline">
		{row.original.poolName}
	</Badge>
{/snippet}

{#snippet usage(row: Row<Subvolume>)}
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

{#snippet createTime(row: Row<Subvolume>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}
