<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		name,
		type,
		namespace,
		health,
		service,
		pod,
		replica,
		container,
		volume,
		nodeport
	};
</script>

{#snippet _row_picker(row: Row<Application>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet name(row: Row<Application>)}
	{row.original.name}
{/snippet}

{#snippet type(row: Row<Application>)}
	<Badge variant="outline">
		{row.original.type}
	</Badge>
{/snippet}

{#snippet namespace(row: Row<Application>)}
	<Badge variant="outline">
		{row.original.namespace}
	</Badge>
{/snippet}

{#snippet health(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.healthies}
	</span>
{/snippet}

{#snippet service(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.services.length}
	</span>
{/snippet}

{#snippet pod(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.pods.length}
	</span>
{/snippet}

{#snippet replica(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.replicas}
	</span>
{/snippet}

{#snippet container(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.containers.length}
	</span>
{/snippet}

{#snippet volume(row: Row<Application>)}
	<span class="flex justify-end">
		{row.original.persistentVolumeClaims.length}
	</span>
{/snippet}

{#snippet nodeport(row: Row<Application>)}
	{#each row.original.services as service}
		{#each service.ports as port}
			{#if port.nodePort > 0}
				<Button
					variant="ghost"
					target="_blank"
					href={`http://${row.original.publicAddress}:${port.nodePort}`}
				>
					{port.targetPort}
					<Icon icon="ph:arrow-square-out" />
				</Button>
			{/if}
		{/each}
	{/each}
{/snippet}
