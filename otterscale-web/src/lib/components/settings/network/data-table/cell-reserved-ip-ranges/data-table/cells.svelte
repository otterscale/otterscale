<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		row_picker,
		type,
		startIp,
		endIp,
		comment
	};
</script>

{#snippet row_picker(row: Row<Network_IPRange>)}
	<Cells.RowPicker {row} />
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
