<script lang="ts" module>
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		name: m.name(),
		state: m.state(),
		in: m.osd_ins(),
		up: m.osd_ups(),
		exists: m.osd_exists(),
		deviceClass: m.device_class(),
		machine: m.machine(),
		placementGroupCount: m.placement_group(),
		usage: m.usage(),
		iops: m.iops()
	};

	export const headers = {
		row_picker,
		name,
		state,
		osdUp,
		osdIn,
		exists,
		deviceClass,
		machine,
		placementGroupCount,
		usage,
		iops,
		actions
	};
</script>

{#snippet row_picker(table: Table<OSD>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.name}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet state(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.state}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet osdUp(column: Column<OSD>)}
	<Layout.Header class="hidden" />
{/snippet}

{#snippet osdIn(column: Column<OSD>)}
	<Layout.Header class="hidden" />
{/snippet}

{#snippet exists(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.exists}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet deviceClass(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.deviceClass}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet machine(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.machine}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet placementGroupCount(column: Column<OSD>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.placementGroupCount}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<OSD>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.usage}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iops(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.iops}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<OSD>)}{/snippet}
