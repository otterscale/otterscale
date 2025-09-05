<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Subvolume_Snapshot } from '$lib/api/storage/v1/storage_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		createTime,
		hasPendingClones,
		actions,
	};
</script>

{#snippet row_picker(table: Table<Subvolume_Snapshot>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Subvolume_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet hasPendingClones(column: Column<Subvolume_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.pending_clones()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet createTime(column: Column<Subvolume_Snapshot>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.create_time()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
