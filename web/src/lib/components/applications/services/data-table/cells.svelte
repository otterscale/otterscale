<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import type { Service } from '../types';

	import CopyButton from '$lib/components/custom/copy-button/copy-button.svelte';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Table from '$lib/components/custom/table';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { m } from '$lib/paraglide/messages';

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
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Service>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Service>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet clusterIp(row: Row<Service>)}
	<Layout.Cell class="items-start">
		{row.original.clusterIp}
	</Layout.Cell>
{/snippet}

{#snippet ports(row: Row<Service>)}
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
							<Table.Head class="text-start">{m.protocol()}</Table.Head>
							<Table.Head class="text-end">{m.port()}</Table.Head>
							<Table.Head class="text-end">{m.node_port()}</Table.Head>
							<Table.Head class="text-end">{m.target_port()}</Table.Head>
							<Table.Head class="text-start">{m.name()}</Table.Head>
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

{#snippet endpoints(row: Row<Service>)}
	{#if row.original.type === 'NodePort'}
		<Layout.Cell class="items-start">
			{#each row.original.ports as port}
				{@const url = `http://${row.original.publicAddress}:${port.nodePort}`}
				<div class="group flex items-center gap-1">
					<Tooltip.Provider>
						<Tooltip.Root delayDuration={13}>
							<Tooltip.Trigger>
								<Badge variant="outline">
									{port.name}
								</Badge>
							</Tooltip.Trigger>
							<Tooltip.Content>
								{url}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
					<CopyButton class="invisible size-4 group-hover:visible" text={url} />
				</div>
			{/each}
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet actions()}{/snippet}
