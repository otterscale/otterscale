<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatTimeAgo } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import { type Subvolume } from './types';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		dataPool: dataPool,
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
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet dataPool(row: Row<Subvolume>)}
	<Badge variant="outline">
		{row.original.dataPool}
	</Badge>
{/snippet}

{#snippet usage(row: Row<Subvolume>)}
	{row.original.usage}
{/snippet}

{#snippet path(row: Row<Subvolume>)}
	<p>
		{row.original.path}
	</p>
{/snippet}

{#snippet mode(row: Row<Subvolume>)}
	<Badge variant="outline">
		{row.original.mode}
	</Badge>
{/snippet}

{#snippet createTime(row: Row<Subvolume>)}
	<p>
		{formatTimeAgo(row.original.createTime)}
	</p>
{/snippet}
