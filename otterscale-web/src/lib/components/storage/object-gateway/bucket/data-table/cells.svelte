<script lang="ts" module>
	import type { Bucket } from '$lib/api/storage/v1/storage_pb';
	import { Cell as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';
	import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		owner,
		usage,
		createTime,
		actions
	};
</script>

{#snippet row_picker(row: Row<Bucket>)}
	<RowPicker {row} />
{/snippet}

{#snippet name(row: Row<Bucket>)}
	{row.original.name}
{/snippet}

{#snippet owner(row: Row<Bucket>)}
	<Badge variant="outline">{row.original.owner}</Badge>
{/snippet}

{#snippet usage(row: Row<Bucket>)}
	{@const { value, unit } = formatCapacity(row.original.usedBytes}
	<div class="flex flex-col items-end">
		<div class="flex items-end">
			{value}
			{unit}
		</div>
		<p class="font-light">{row.original.usedObjects} unit(s)</p>
	</div>
{/snippet}

{#snippet createTime(row: Row<Bucket>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}

{#snippet actions(row: Row<Bucket>)}
	<Actions bucket={row.original} />
{/snippet}
