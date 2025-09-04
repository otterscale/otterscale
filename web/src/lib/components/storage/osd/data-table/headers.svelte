<script lang="ts" module>
	import type { OSD } from '$lib/api/storage/v1/storage_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

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
		actions,
	};
</script>

{#snippet row_picker(table: Table<OSD>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet state()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.state()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet osdUp()}
	<Layout.Header class="hidden" />
{/snippet}

{#snippet osdIn()}
	<Layout.Header class="hidden" />
{/snippet}

{#snippet exists(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.osd_exists()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet deviceClass(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.device_class()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet machine(column: Column<OSD>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.machine()}</Layout.HeaderViewer>
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
		<Layout.HeaderViewer>{m.placement_group()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<OSD>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.usage()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet iops()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.iops()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
