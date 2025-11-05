<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import { type LargeLanguageModel } from '../type';

	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_expander,
		row_picker,
		name,
		namespace,
		chart_version,
		application_version,
		status,
		description,
		requests,
		limits,
		first_deployed_at,
		last_deployed_at,
		relation,
		pods,
		action,
	};
</script>

{#snippet row_expander()}{/snippet}

{#snippet row_picker(table: Table<LargeLanguageModel>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet chart_version(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.version()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet application_version(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.version()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.status()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet description(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.description()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet requests(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.requests()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet limits(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.limits()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet first_deployed_at(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.time()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet last_deployed_at(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.time()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet pods(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.pods()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet relation()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.relational()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet action()}{/snippet}
