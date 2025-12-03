<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Repository } from '$lib/api/registry/v1/registry_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		manifests,
		sizeBytes,
		latestTag,
		actions
	};
</script>

{#snippet row_picker(table: Table<Repository>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Repository>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet manifests(column: Column<Repository>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.manifests()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet sizeBytes(column: Column<Repository>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.size()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet latestTag()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.latest_tag()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet actions()}{/snippet}
