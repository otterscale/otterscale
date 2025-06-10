<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	import type { Row } from '@tanstack/table-core';
	import type { BlockImage } from './types';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		pool: pool,
		namespace: namespace,
		size: size,
		usage: usage,
		objects: objects,
		objectSize: objectSize,
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

{#snippet size(row: Row<BlockImage>)}
	{row.original.size}
{/snippet}

{#snippet usage(row: Row<BlockImage>)}
	{row.original.usage}
{/snippet}

{#snippet objects(row: Row<BlockImage>)}
	{row.original.objects}
{/snippet}

{#snippet objectSize(row: Row<BlockImage>)}
	{row.original.objectSize}
{/snippet}

{#snippet parent(row: Row<BlockImage>)}
	<Badge variant="outline">
		{row.original.parent}
	</Badge>
{/snippet}

{#snippet mirroring(row: Row<BlockImage>)}
	<Badge variant="outline">
		{row.original.mirroring}
	</Badge>
{/snippet}

{#snippet nextScheduledSnapshot(row: Row<BlockImage>)}
	<Badge variant="outline">
		{row.original.nextScheduledSnapshot}
	</Badge>
{/snippet}
