<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';

	import Actions from './actions.svelte';

	export const cells = {
		row_picker,
		type,
		startIp,
		endIp,
		comment,
		actions
	};
</script>

{#snippet row_picker(row: Row<Network_IPRange>)}
	<Table.Cell alignClass="items-start">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet type(row: Row<Network_IPRange>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet startIp(row: Row<Network_IPRange>)}
	<Table.Cell alignClass="items-start">
		{row.original.startIp}
	</Table.Cell>
{/snippet}

{#snippet endIp(row: Row<Network_IPRange>)}
	<Table.Cell alignClass="items-start">
		{row.original.endIp}
	</Table.Cell>
{/snippet}

{#snippet comment(row: Row<Network_IPRange>)}
	<Table.Cell alignClass="items-start">
		{row.original.comment}
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Network_IPRange>; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions ipRange={data.row.original} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
