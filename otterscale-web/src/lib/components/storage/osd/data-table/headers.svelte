<script lang="ts" module>
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import * as Layout from '$lib/components/custom/data-table/data-table-layout';
	import { Header as RowPicker } from '$lib/components/custom/data-table/data-table-row-pickers';
	import Sorter from '$lib/components/custom/data-table/data-table-sorter.svelte';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		name: m.storage_osds_heads_name(),
		state: m.storage_osds_heads_state(),
		in: m.storage_osds_heads_osd_ins(),
		up: m.storage_osds_heads_osd_ups(),
		exists: m.storage_osds_heads_exists(),
		deviceClass: m.storage_osds_heads_device_class(),
		machine: m.storage_osds_heads_machine(),
		placementGroupCount: m.storage_osds_heads_placement_group_count(),
		usage: m.storage_osds_heads_usage(),
		iops: m.storage_osds_heads_iops()
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
