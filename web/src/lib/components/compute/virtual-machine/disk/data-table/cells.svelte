<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import {
		VirtualMachine_Disk_Bus,
		VirtualMachine_Disk_Volume_Source_Type
	} from '$lib/api/instance/v1/instance_pb';
	import type { EnhancedDisk } from '$lib/components/compute/virtual-machine/units/type';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import * as Tooltip from '$lib/components/ui/tooltip';
	import { formatCapacity } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		bus,
		bootOrder,
		dataVolume,
		type,
		phase,
		boot,
		size,
		actions
	};
</script>

{#snippet row_picker(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet bus(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{VirtualMachine_Disk_Bus[row.original.bus]}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet bootOrder(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{#if row.original.bootOrder}
			<Badge variant="outline">
				{row.original.bootOrder}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet dataVolume(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{row.original.volume?.name ?? ''}
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{#if row.original.volume?.source?.type}
			<Badge variant="outline">
				{VirtualMachine_Disk_Volume_Source_Type[row.original.volume.source.type]}
			</Badge>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet phase(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					{#if row.original.phase === 'Succeeded'}
						<Icon icon="ph:check" class="text-green-600" />
					{:else if row.original.phase === 'ImportInProgress'}
						<Icon icon="ph:spinner" class="animate-spin text-blue-600" />
					{:else}
						<Icon icon="ph:x" class="text-gray-400" />
					{/if}
				</Tooltip.Trigger>
				<Tooltip.Content>
					{row.original.phase}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	</Layout.Cell>
{/snippet}

{#snippet boot(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		{#if row.original.bootImage}
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{#if row.original.bootImage === true}
							<Icon icon="ph:power" class="text-green-600" />
						{:else}
							<Icon icon="ph:power" class="text-gray-400" />
						{/if}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{row.original.bootImage}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet size(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-end">
		{#if row.original.sizeBytes}
			{@const { value, unit } = formatCapacity(row.original.sizeBytes)}
			{value}
			{unit}
		{/if}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<EnhancedDisk>)}
	<Layout.Cell class="items-start">
		<Actions enhancedDisk={row.original} />
	</Layout.Cell>
{/snippet}
