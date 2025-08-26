<script lang="ts" module>
	import * as Layout from '$lib/components/custom/data-table/layout';
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
