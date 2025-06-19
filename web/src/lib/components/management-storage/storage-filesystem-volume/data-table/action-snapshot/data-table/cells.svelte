<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	import type { Row } from '@tanstack/table-core';
	import { type Snapshot } from './types';
	import { formatTimeAgo } from '$lib/formatter';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		pendingClones: pendingClones
	};
</script>

{#snippet _row_picker(row: Row<Snapshot>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Snapshot>)}
	<p>
		{row.original.name}
	</p>
{/snippet}

{#snippet createTime(row: Row<Snapshot>)}
	<p>
		{formatTimeAgo(row.original.createTime)}
	</p>
{/snippet}

{#snippet pendingClones(row: Row<Snapshot>)}
	<span class="flex items-center gap-2">
		{#each row.original.pendingClones as pendingClone}
			<Badge variant="outline">
				{pendingClone}
			</Badge>
		{/each}
	</span>
{/snippet}
