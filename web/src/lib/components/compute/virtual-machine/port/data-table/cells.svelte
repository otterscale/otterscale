<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Application_Service_Port } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		row_picker,
		name,
		protocol,
		port,
		nodePort,
	};
</script>

{#snippet row_picker(row: Row<Application_Service_Port>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Application_Service_Port>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet protocol(row: Row<Application_Service_Port>)}
	<Layout.Cell class="items-start">
		{#if row.original.protocol}
			<Badge variant="outline">
				{row.original.protocol}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet port(row: Row<Application_Service_Port>)}
	<Layout.Cell class="items-start">
		{#if row.original.port}
			<Badge variant="outline">
				{row.original.port}:{row.original.nodePort}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet nodePort(row: Row<Application_Service_Port>)}
	<Layout.Cell class="items-start">
		{#if row.original.nodePort}
			<Badge variant="outline">
				{row.original.nodePort}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}
