<script lang="ts" module>
	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		type,
		startIp,
		endIp,
		comment
	};
</script>

{#snippet _row_picker(row: Row<Network_IPRange>)}
	<TableRowPicker {row} />
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
