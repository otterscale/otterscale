<script lang="ts" module>
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatCapacity } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import type { BlockImage } from './types';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		pool: pool,
		namespace: namespace,
		usage: usage,
		objects: objects,
		parent: parent,
		mirroring: mirroring,
		nextScheduledSnapshot: nextScheduledSnapshot
	};
</script>

{#snippet _row_picker(row: Row<BlockImage>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<BlockImage>)}
	{row.original.name}
{/snippet}

{#snippet pool(row: Row<BlockImage>)}
	{row.original.pool}
{/snippet}

{#snippet namespace(row: Row<BlockImage>)}
	<Badge variant="outline">
		{row.original.namespace}
	</Badge>
{/snippet}

{#snippet usage(row: Row<BlockImage>)}
	{@const total = row.original.size}
	{@const used = Math.round((row.original.size * row.original.usage) / 100)}
	<Progress.Root numerator={used} denominator={total}>
		{#snippet ratio({ numerator, denominator })}
			{Math.round((numerator * 100) / denominator)}%
		{/snippet}
		{#snippet detail({ numerator, denominator })}
			{numerator}/{denominator}
		{/snippet}
	</Progress.Root>
{/snippet}

{#snippet objects(row: Row<BlockImage>)}
	{@const objectSize = formatCapacity(row.original.objectSize)}
	<div class="flex flex-col items-end">
		{row.original.objects}
		<p class="text-muted-foreground font-light">{objectSize.value} {objectSize.unit}</p>
	</div>
{/snippet}

{#snippet parent(row: Row<BlockImage>)}
	{#if row.original.parent}
		<Badge variant="outline">
			{row.original.parent}
		</Badge>
	{:else}
		<Badge variant="secondary">Null</Badge>
	{/if}
{/snippet}

{#snippet mirroring(row: Row<BlockImage>)}
	<Badge variant="outline">
		{row.original.mirroring}
	</Badge>
{/snippet}

{#snippet nextScheduledSnapshot(row: Row<BlockImage>)}
	{row.original.nextScheduledSnapshot}
{/snippet}
