<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { RowPickers } from '$lib/components/custom/data-table';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker,
		type,
		startIp,
		endIp,
		comment
	};
</script>

{#snippet _row_picker(row: Row<Network_IPRange>)}
	<RowPickers.Cell {row} />
{/snippet}

{#snippet type(row: Row<Network_IPRange>)}
	<Badge variant="outline">
		{row.original.type}
	</Badge>
{/snippet}

{#snippet startIp(row: Row<Network_IPRange>)}
	{row.original.startIp}
{/snippet}

{#snippet endIp(row: Row<Network_IPRange>)}
	{row.original.endIp}
{/snippet}

{#snippet comment(row: Row<Network_IPRange>)}
	{row.original.comment}
{/snippet}
