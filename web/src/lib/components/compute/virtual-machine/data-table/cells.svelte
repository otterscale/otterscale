<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import Actions from './cell-actions.svelte';

	import { page } from '$app/state';
	import type { VirtualMachine } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { Disk } from '$lib/components/compute/virtual-machine/disk';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Badge } from '$lib/components/ui/badge';
	import { formatCapacity } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		namespace,
		network,
		node,
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
		<a
			class="underline hover:no-underline"
			href={`${page.url}/${row.original.metadata?.namespace}/${row.original.metadata?.name}`}
		>
			{row.original.metadata?.name}
		</a>
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
