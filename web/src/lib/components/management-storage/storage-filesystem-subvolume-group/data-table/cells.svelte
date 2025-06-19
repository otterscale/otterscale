<script lang="ts" module>
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { formatTimeAgo } from '$lib/formatter';
	import type { Row } from '@tanstack/table-core';
	import { type SubvolumeGroup } from './types';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		dataPool: dataPool,
		usage: usage,
		mode: mode,
		createTime: createTime
	};
</script>

{#snippet _row_picker(row: Row<SubvolumeGroup>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<SubvolumeGroup>)}
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet dataPool(row: Row<SubvolumeGroup>)}
	<Badge variant="outline">
		{row.original.dataPool}
	</Badge>
{/snippet}

{#snippet usage(row: Row<SubvolumeGroup>)}
	{row.original.usage}
{/snippet}

{#snippet mode(row: Row<SubvolumeGroup>)}
	<Badge variant="outline">
		{row.original.mode}
	</Badge>
{/snippet}

{#snippet createTime(row: Row<SubvolumeGroup>)}
	<p>
		{formatTimeAgo(row.original.createTime)}
	</p>
{/snippet}
