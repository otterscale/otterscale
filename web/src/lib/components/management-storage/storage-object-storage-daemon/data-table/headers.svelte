<script lang="ts" module>
	import Sorter from '$lib/components/custom/data-table/data-table-column-sorter.svelte';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Column, Table } from '@tanstack/table-core';
	import type { ObjectStorageDaemon } from './types';

	export const headers = {
		_row_picker: _row_picker,
		id: id,
		host: host,
		devices: devices,
		status: status,
		deviceClass: deviceClass,
		pgs: pgs,
		size: size,
		flags: flags,
		usage: usage,
		readBytes: readBytes,
		writeBytes: writeBytes
	};
</script>

{#snippet _row_picker(table: Table<ObjectStorageDaemon>)}
	<Checkbox
		checked={table.getIsAllPageRowsSelected()}
		indeterminate={table.getIsSomePageRowsSelected() && !table.getIsAllPageRowsSelected()}
		onCheckedChange={(value) => table.toggleAllPageRowsSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select all"
	/>
{/snippet}
{#snippet id(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>ID</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet host(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>HOST</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet devices(column: Column<ObjectStorageDaemon>)}
	<Layout.Header class="flex justify-end">
		<Layout.HeaderViewer>DEVICES</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>STATUS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet deviceClass(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>DEVICE CLASS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet pgs(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>PLACEMENT GROUP</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet size(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>SIZE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet flags(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>FLAG</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>USAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet readBytes(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>READ</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet writeBytes(column: Column<ObjectStorageDaemon>)}
	<Layout.Header>
		<Layout.HeaderViewer>WRITE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
