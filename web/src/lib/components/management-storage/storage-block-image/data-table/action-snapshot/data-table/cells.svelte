<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { BlockImageSnapshot } from './types';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity, formatTimeAgo } from '$lib/formatter';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		size: size,
		used: used,
		state: state,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<BlockImageSnapshot>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<BlockImageSnapshot>)}
	{row.original.name}
{/snippet}

{#snippet size(row: Row<BlockImageSnapshot>)}
	{@const size = formatCapacity(row.original.size)}
	<div class="flex flex-col items-end">
		{row.original.size}
		<p class="text-muted-foreground font-light">{size.value} {size.unit}</p>
	</div>
{/snippet}

{#snippet used(row: Row<BlockImageSnapshot>)}
	{@const used = formatCapacity(row.original.used)}
	<div class="flex flex-col items-end">
		{row.original.used}
		<p class="text-muted-foreground font-light">{used.value} {used.unit}</p>
	</div>
{/snippet}

{#snippet state(row: Row<BlockImageSnapshot>)}
	<Badge variant="outline">
		{row.original.state}
	</Badge>
{/snippet}

{#snippet createTime(row: Row<BlockImageSnapshot>)}
	{formatTimeAgo(new Date(row.original.createTime))}
{/snippet}
