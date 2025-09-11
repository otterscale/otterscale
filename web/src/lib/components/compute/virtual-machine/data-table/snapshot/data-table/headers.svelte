<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { VirtualMachineSnapshot } from '$lib/api/kubevirt/v1/kubevirt_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		namespace,
		description,
		statusPhase,
		lastConditionMessage,
		lastConditionReason,
	};
</script>

{#snippet row_picker(table: Table<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet statusPhase(column: Column<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>Status Phase</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet description(column: Column<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>Description</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet lastConditionMessage(column: Column<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>Condition Message</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet lastConditionReason(column: Column<VirtualMachineSnapshot>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>Condition Reason</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

<!-- {#snippet actions()}{/snippet} -->
