<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Application_Service_Port } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		row_picker,
		name,
		protocol,
		port,
		nodePort
	};
</script>

{#snippet row_picker(row: Row<Application_Service_Port>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Application_Service_Port>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet protocol(row: Row<Application_Service_Port>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.protocol}
			<Badge variant="outline">
				{row.original.protocol}
			</Badge>
		{/if}
	</Table.Cell>
{/snippet}

{#snippet port(row: Row<Application_Service_Port>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.port}
			<Badge variant="outline">
				{row.original.port}:{row.original.nodePort}
			</Badge>
		{/if}
	</Table.Cell>
{/snippet}

{#snippet nodePort(row: Row<Application_Service_Port>)}
	<Table.Cell alignClass="items-start">
		{#if row.original.nodePort}
			<Badge variant="outline">
				{row.original.nodePort}
			</Badge>
		{/if}
	</Table.Cell>
{/snippet}
