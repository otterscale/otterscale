<script lang="ts" module>
	import type { Machine } from '$lib/api/machine/v1/machine_pb';
	import { Layout, RowPickers, Sorter } from '$lib/components/custom/data-table';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		_row_picker,
		fqdn_ip,
		powerState,
		status,
		cores_arch,
		ram,
		disk,
		storage,
		tags,
		actions
	};
</script>

{#snippet _row_picker(table: Table<Machine>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPickers.Header {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet fqdn_ip(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>
			<p>FQDN</p>
			<p>IP</p>
		</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet powerState(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>POWER</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<Machine>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>STATUS</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet cores_arch(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>
			<p>CORES</p>
			<p>ARCH</p>
		</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ram(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>RAM</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet disk(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>DISK</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet storage(column: Column<Machine>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>STORAGE</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet tags()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>TAGS</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
