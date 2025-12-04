<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Table from '$lib/components/custom/table/index.js';
	import { TagGroup } from '$lib/components/tag-group';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import { m } from '$lib/paraglide/messages';

	import type { Service } from '../types';

	export const cells = {
		row_picker,
		name,
		type,
		clusterIp,
		ports,
		endpoints,
		actions
	};
</script>

{#snippet row_picker(row: Row<Service>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Service>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet type(row: Row<Service>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet clusterIp(row: Row<Service>)}
	<Table.Cell alignClass="items-start">
		{row.original.clusterIp}
	</Table.Cell>
{/snippet}

{#snippet ports(row: Row<Service>)}
	<Table.Cell alignClass="items-end">
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
							<Table.Head class="text-start">{m.protocol()}</Table.Head>
							<Table.Head class="text-end">{m.port()}</Table.Head>
							<Table.Head class="text-end">{m.node_port()}</Table.Head>
							<Table.Head class="text-end">{m.target_port()}</Table.Head>
							<Table.Head class="text-start">{m.name()}</Table.Head>
						</Table.Row>
					</Table.Header>
					<Table.Body>
						{#each row.original.ports as port, index (index)}
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
	</Table.Cell>
{/snippet}

{#snippet endpoints(row: Row<Service>)}
	{#if row.original.type === 'NodePort'}
		<Table.Cell alignClass="items-start">
			<TagGroup
				items={row.original.ports.map((port) => ({
					title: port.name ?? '',
					description: `http://${row.original.hostname}:${port.nodePort}`,
					icon: 'ph:tag'
				}))}
			/>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet actions()}{/snippet}
