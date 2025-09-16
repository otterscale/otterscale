<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import { page } from '$app/state';
	import { type VirtualMachine, VirtualMachine_status } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { Disk } from '$lib/components/compute/virtual-machine/disk';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';

	function getStatusIcon(status: VirtualMachine_status) {
		switch (status) {
			case VirtualMachine_status.RUNNING:
				return { icon: 'ph:power', color: 'text-green-600' };
			case VirtualMachine_status.STOPPED:
				return { icon: 'ph:power', color: 'text-gray-600' };
			case VirtualMachine_status.PAUSED:
				return { icon: 'ph:pause-circle', color: 'text-yellow-600' };
			case VirtualMachine_status.STARTING:
				return { icon: 'ph:arrow-clockwise', color: 'text-blue-500 animate-spin' };
			case VirtualMachine_status.PROVISIONING:
				return { icon: 'ph:gear', color: 'text-blue-500 animate-spin' };
			case VirtualMachine_status.TERMINATING:
				return { icon: 'ph:trash', color: 'text-red-500' };
			case VirtualMachine_status.UNKNOWN:
			default:
				return { icon: 'ph:warning-circle-fill', color: 'text-amber-500' };
		}
	}

	export const cells = {
		row_picker,
		name,
		namespace,
		network,
		node,
		status,
		cpu,
		memory,
		disk,
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
		{row.original.metadata?.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.metadata?.namespace}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet network(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.networkName}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet node(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		{row.original.nodeName}
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<VirtualMachine>)}
	{@const statusIcon = getStatusIcon(row.original.statusPhase)}
	<Layout.Cell class="items-start">
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<Icon icon={statusIcon.icon} class={`${statusIcon.color} h-5 w-5`} />
				</Tooltip.Trigger>
				<Tooltip.Content>
					{VirtualMachine_status[row.original.statusPhase]}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Layout.Cell>
{/snippet}

{#snippet cpu(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-end">
		{row.original.resources?.cpuCores}
	</Layout.Cell>
{/snippet}

{#snippet memory(row: Row<VirtualMachine>)}
	{@const memory = Number(row.original.resources?.memoryBytes)}
	<Layout.Cell class="items-end">
		{@const { value: memoryValue, unit: memoryUnit } = formatCapacity(memory)}
		{memoryValue}
		{memoryUnit}
	</Layout.Cell>
{/snippet}

{#snippet disk(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-end">
		<Disk virtualMachines={row.original.disks} />
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<VirtualMachine>)}
	<Layout.Cell class="items-start">
		<Actions virtualMachine={row.original} />
	</Layout.Cell>
{/snippet}
