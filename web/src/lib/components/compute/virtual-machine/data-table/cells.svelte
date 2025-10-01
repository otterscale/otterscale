<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import { page } from '$app/state';
	import type { VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import { Disk } from '$lib/components/compute/virtual-machine/disk';
	import { Port } from '$lib/components/compute/virtual-machine/port';
	import { getStatusInfo } from '$lib/components/compute/virtual-machine/units/type';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Table from '$lib/components/ui/table';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';
	import { dynamicPaths } from '$lib/path';

	export const cells = {
		row_picker,
		name,
		status,
		namespace,
		machineId,
		instanceType,
		disk,
		port,
		createTime,
		actions,
	};
</script>

{#snippet row_picker(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<div class="flex items-center gap-1">
			{row.original.name}
		</div>
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<VirtualMachine>)}
	{@const statusInfo = getStatusInfo(row.original.status)}
	<Layout.Cell class="items-start">
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Icon icon={statusInfo.icon} class={`${statusInfo.color} h-5 w-5`} />
				</Tooltip.Trigger>
				<Tooltip.Content>
					{statusInfo.text}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.namespace}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet machineId(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.machineId}
			<a
				class="m-0 p-0 underline hover:no-underline"
				href={`${dynamicPaths.machinesMetal(page.params.scope).url}/${row.original.machineId}`}
			>
				{row.original.hostname}
			</a>
			<Layout.SubCell>
				<span class="text-muted-foreground flex items-center gap-1 text-xs">
					{row.original.ipAddresses}
				</span>
			</Layout.SubCell>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet instanceType(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.instanceType}
			<div class="flex items-center gap-1">
				<Badge variant="outline">
					{row.original.instanceType.name}
				</Badge>
				<HoverCard.Root>
					<HoverCard.Trigger>
						<Icon icon="ph:info" />
					</HoverCard.Trigger>
					<HoverCard.Content class="min-w-[300px]">
						<Table.Root>
							<Table.Body class="text-xs">
								{#if row.original.instanceType.name}
									<Table.Row>
										<Table.Head class="text-left">Name</Table.Head>
										<Table.Cell>
											<Badge variant="outline">{row.original.instanceType.name}</Badge>
										</Table.Cell>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.namespace}
									<Table.Row>
										<Table.Head class="text-left">Namespace</Table.Head>
										<Table.Cell>{row.original.instanceType.namespace}</Table.Cell>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.cpuCores}
									<Table.Row>
										<Table.Head class="text-left">CPU Cores</Table.Head>
										<Table.Cell>{row.original.instanceType.cpuCores}</Table.Cell>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.memoryBytes}
									<Table.Row>
										<Table.Head class="text-left">Memory</Table.Head>
										<Table.Cell
											>{Number(row.original.instanceType.memoryBytes) / 1024 ** 3} GB</Table.Cell
										>
									</Table.Row>
								{/if}
								{#if row.original.instanceType.clusterWide !== undefined}
									<Table.Row>
										<Table.Head class="text-left">Cluster Wide</Table.Head>
										<Table.Cell>
											<Badge variant="outline">{row.original.instanceType.clusterWide}</Badge>
										</Table.Cell>
									</Table.Row>
								{/if}
							</Table.Body>
						</Table.Root>
					</HoverCard.Content>
				</HoverCard.Root>
			</div>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet disk(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-end">
		<Disk virtualMachine={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet port(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-end">
		<Port virtualMachine={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet createTime(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.createdAt}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.createdAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.createdAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<Actions virtualMachine={row.original} />
	</Layout.Cell>
{/snippet}
