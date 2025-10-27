<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import { type LargeLanguageModel } from '../type';

	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		model,
		name,
		replicas,
		healthies,
		gpu_cache,
		kv_cache,
		requests,
		time_to_first_token,
		relation,
		action,
	};
</script>

{#snippet row_picker(table: Table<LargeLanguageModel>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet model()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.model()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.model_name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet replicas(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.replica()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet healthies(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.health()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet gpu_cache(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.gpu_cache()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet kv_cache(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.kv_cache()}</Layout.HeaderViewer>
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

{#snippet time_to_first_token(column: Column<LargeLanguageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.uptime()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet relation()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.gpu_relation()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet action()}{/snippet}
