<script lang="ts" module>
	import type { User_Key } from '$lib/api/storage/v1/storage_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const headers = {
		row_picker,
		accessKey,
		actions,
	};
</script>

{#snippet row_picker(table: Table<User_Key>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet accessKey(column: Column<User_Key>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.access_key()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
