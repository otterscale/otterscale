<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Network_IPRange } from '$lib/api/network/v1/network_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import Actions from './actions.svelte';

	export const cells = {
		row_picker,
		type,
		startIp,
		endIp,
		comment,
		actions,
	};
</script>

{#snippet row_picker(row: Row<Network_IPRange>)}
	<Layout.Cell class="items-start">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Network_IPRange>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet startIp(row: Row<Network_IPRange>)}
	<Layout.Cell class="items-start">
		{row.original.startIp}
	</Layout.Cell>
{/snippet}

{#snippet endIp(row: Row<Network_IPRange>)}
	<Layout.Cell class="items-start">
		{row.original.endIp}
	</Layout.Cell>
{/snippet}

{#snippet comment(row: Row<Network_IPRange>)}
	<Layout.Cell class="items-start">
		{row.original.comment}
	</Layout.Cell>
{/snippet}


{#snippet actions(data: {row: Row<Network_IPRange>; reloadManager: ReloadManager})}
	<Layout.Cell class="items-start">
		<Actions ipRange={data.row.original} reloadManager={data.reloadManager} />
	</Layout.Cell>
{/snippet}