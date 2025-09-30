<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import type { VirtualMachine } from '$lib/api/virtual_machine/v1/virtual_machine_pb';
	import { Disk } from '$lib/components/compute/virtual-machine/disk';
	import { Port } from '$lib/components/compute/virtual-machine/port';
	import { getStatusInfo, getInstancePhaseInfo } from '$lib/components/compute/virtual-machine/units/type';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatTimeAgo } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		status,
		namespace,
		machineId,
		instanceTypeName,
		instancePhase,
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
			{#if !row.original.ready}
				<Icon icon="ph:spinner-gap" class="size-5 animate-spin" />
			{/if}
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
			<Badge variant="outline">
				{row.original.machineId}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet instanceTypeName(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{#if row.original.instanceTypeName}
			<Badge variant="outline">
				{row.original.instanceTypeName}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet instancePhase(row: Row<VirtualMachine>)}
	{@const instancePhaseInfo = getInstancePhaseInfo(row.original.instancePhase)}
	<Layout.Cell class="items-start">
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Icon icon={instancePhaseInfo.icon} class={`${instancePhaseInfo.color} h-5 w-5`} />
				</Tooltip.Trigger>
				<Tooltip.Content>
					{instancePhaseInfo.text}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
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
