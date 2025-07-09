<script lang="ts" module>
	import type { MON } from '$gen/api/storage/v1/storage_pb';
	import { Badge } from '$lib/components/ui/badge';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		leader: leader,
		name: name,
		rank: rank,
		publicAddress: publicAddress
	};
</script>

{#snippet _row_picker(row: Row<MON>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet leader(row: Row<MON>)}
	<Badge variant="outline">
		{row.original.leader}
	</Badge>
{/snippet}

{#snippet name(row: Row<MON>)}
	{row.original.name}
{/snippet}

{#snippet rank(row: Row<MON>)}
	<div class="text-end">{row.original.rank}</div>
{/snippet}

{#snippet publicAddress(row: Row<MON>)}
	<Badge variant="outline">
		{row.original.publicAddress}
	</Badge>
{/snippet}
