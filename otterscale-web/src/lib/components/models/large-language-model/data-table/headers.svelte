<script lang="ts" module>
	import { Layout, RowPickers, Sorter } from '$lib/components/custom/data-table';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';
	import { type LargeLangeageModel } from '../protobuf.svelte';

	export const messages = {
		name: m.models_large_language_model_headers_name(),
		version: m.models_large_language_model_headers_version(),
		parameters: m.models_large_language_model_headers_parameters(),
		accuracy: m.models_large_language_model_headers_accuracy(),
		speed: m.models_large_language_model_headers_speed(),
		architecture: m.models_large_language_model_headers_architecture(),
		requests: m.models_large_language_model_headers_requests(),
		uptime: m.models_large_language_model_headers_uptime()
	};

	export const headers = {
		row_picker,
		name,
		version,
		parameters,
		accuracy,
		speed,
		architecture,
		requests,
		uptime
	};
</script>

{#snippet row_picker(table: Table<LargeLangeageModel>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<RowPickers.Header {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages['name']}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}
{#snippet version(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages['version']}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet parameters(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['parameters']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet accuracy(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['accuracy']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet speed(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['speed']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet architecture(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages['architecture']}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet requests(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['requests']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet uptime(column: Column<LargeLangeageModel>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages['uptime']}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
