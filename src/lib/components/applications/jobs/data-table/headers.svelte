<script lang="ts" module>
	import type { Column, Table } from '@tanstack/table-core';

	import type { Job } from '$lib/api/application/v1/application_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';

	export const headers = {
		row_picker,
		name,
		namespace,
		active,
		ready,
		succeeded,
		failed,
		terminating,
		status,
		conditions,
		startedAt,
		completedAt
	};
</script>

{#snippet row_picker(table: Table<Job>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Job>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.name()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<Job>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.namespace()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet active()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.active()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet ready()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.ready()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet terminating()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.terminating()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet succeeded()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.succeeded()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet failed()}
	<Layout.Header class="justify-end">
		<Layout.HeaderViewer>{m.failed()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet status(column: Column<Job>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.status()}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet conditions()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{m.conditions()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet startedAt(column: Column<Job>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.started_at()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet completedAt(column: Column<Job>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{m.completed_at()}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
