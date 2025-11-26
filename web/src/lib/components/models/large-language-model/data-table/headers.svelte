<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_expander,
		row_picker,
		id,
		name,
		namespace,
		status,
		description,
		first_deployed_at,
		last_deployed_at,
		chart_version,
		app_version,
		prefill,
		decode,
		gpu_relation,
		action
	};
</script>

{#snippet row_expander(table: Table<Model>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowExpander {table} />
		</Layout.HeaderController>
	</Layout.Header>{/snippet}

{#snippet row_picker(table: Table<Model>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet id(column: Column<Model>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.model_name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Model>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<Model>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<Model>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.status()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet description(column: Column<Model>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.description()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet first_deployed_at(column: Column<Model>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.first_deployed_at()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet last_deployed_at(column: Column<Model>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.last_deployed_at()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet chart_version(column: Column<Model>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.chart_version()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet app_version(column: Column<Model>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.application_version()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet prefill()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.prefill()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet decode()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.decode()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet gpu_relation()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.gpu_relation()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet action()}{/snippet}
