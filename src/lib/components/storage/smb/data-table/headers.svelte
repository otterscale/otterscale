<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { SMBShare } from '$lib/api/storage/v1/storage_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		uri,
		healthies,
		size,
		browsable,
		read_only,
		guest_ok,
		map_to_guest,
		mode,
		valid_users,
		actions
	};
</script>

{#snippet row_picker(table: Table<SMBShare>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet uri(column: Column<SMBShare>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.uri()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet healthies(column: Column<SMBShare>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.healthies()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet size(column: Column<SMBShare>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.size()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet browsable()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.browsable()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet read_only()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.read_only()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet guest_ok()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.guest_accessible()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet map_to_guest(column: Column<SMBShare>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.map_to_guest()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet mode(column: Column<SMBShare>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.mode()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet valid_users(column: Column<SMBShare>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.valid_users()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
