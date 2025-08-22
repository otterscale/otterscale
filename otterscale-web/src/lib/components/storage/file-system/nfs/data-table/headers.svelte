<script lang="ts" module>
	import type { Subvolume } from '$lib/api/storage/v1/storage_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		name: m.name(),
		poolName: m.pool_name(),
		usage: m.usage(),
		path: m.path(),
		mode: m.mode(),
		createTime: m.create_time(),
		exportSubvolume: m.export(),
		snapshots: m.snapshot()
	};

	export const headers = {
		row_picker,
		name,
		poolName,
		usage,
		path,
		mode,
		createTime,
		exportSubvolume,
		snapshots,
		actions
	};
</script>

{#snippet row_picker(table: Table<Subvolume>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet path(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.path()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet mode(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.mode()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet poolName(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.pool_name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet usage(column: Column<Subvolume>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.usage()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet exportSubvolume(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.export()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Subvolume>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.create_time()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet snapshots(column: Column<Subvolume>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.snapshot()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions(column: Column<Subvolume>)}{/snippet}
