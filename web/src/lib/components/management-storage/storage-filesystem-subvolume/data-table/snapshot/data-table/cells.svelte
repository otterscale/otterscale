<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { Subvolume_Snapshot } from '$gen/api/storage/v1/storage_pb';
	import { formatTimeAgo } from '$lib/formatter';
	import { Badge } from '$lib/components/ui/badge';
	import { timestampDate } from '@bufbuild/protobuf/wkt';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		createTime: createTime,
		hasPendingClones: hasPendingClones
	};
</script>

{#snippet _row_picker(row: Row<Subvolume_Snapshot>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Subvolume_Snapshot>)}
	{row.original.name}
{/snippet}

{#snippet hasPendingClones(row: Row<Subvolume_Snapshot>)}
	<Badge variant="outline">{row.original.hasPendingClones}</Badge>
{/snippet}

{#snippet createTime(row: Row<Subvolume_Snapshot>)}
	{#if row.original.createdAt}
		{formatTimeAgo(timestampDate(row.original.createdAt))}
	{/if}
{/snippet}
