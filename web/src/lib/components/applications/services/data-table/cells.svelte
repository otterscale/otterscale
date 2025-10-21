<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { Application_Service } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card';

	export const cells = {
		row_picker,
		name,
		type,
		clusterIp,
		ports,
		endpoints,
		actions,
	};
</script>

{#snippet row_picker(row: Row<Application_Service>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet clusterIp(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		{row.original.clusterIp}
	</Layout.Cell>
{/snippet}

{#snippet ports(row: Row<Application_Service>)}
	<Layout.Cell class="items-end">
		<HoverCard.Root>
			<HoverCard.Trigger>
				<span class="flex items-center justify-center gap-1">
					{row.original.ports.length}
					<Icon icon="ph:info" />
				</span>
			</HoverCard.Trigger>
			<HoverCard.Content class="m-0 w-fit p-0">
				<Table.Root>
					<Table.Header>
						<Table.Row>
							<Table.Head class="text-start">Protocol</Table.Head>
							<Table.Head class="text-end">Port</Table.Head>
							<Table.Head class="text-end">Node Port</Table.Head>
							<Table.Head class="text-end">Target Port</Table.Head>
							<Table.Head class="text-start">Name</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each row.original.ports as port}
							<Table.Row>
								<Table.Cell class="text-start">
									{port.protocol}
								</Table.Cell>
								<Table.Cell class="text-end">
									{port.port}
								</Table.Cell>
								<Table.Cell class="text-end">
									{port.nodePort}
								</Table.Cell>
								<Table.Cell class="text-end">
									{port.targetPort}
								</Table.Cell>
								<Table.Cell class="text-start">
									{port.name}
								</Table.Cell>
							</Table.Row>
						{/each}
					</Table.Body>
				</Table.Root>
			</HoverCard.Content>
		</HoverCard.Root>
	</Layout.Cell>
{/snippet}

{#snippet endpoints()}{/snippet}

{#snippet actions(row: Row<Application_Service>)}
	<Layout.Cell class="items-start">
		<Actions service={row.original} />
	</Layout.Cell>
{/snippet}
