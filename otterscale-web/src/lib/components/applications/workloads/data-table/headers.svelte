<script lang="ts" module>
	import type { Application } from '$lib/api/application/v1/application_pb';
	import { Headers, Sorter } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { m } from '$lib/paraglide/messages';
	import type { Column, Table } from '@tanstack/table-core';

	export const messages = {
		name: m.applications_workloads_heads_name(),
		type: m.applications_workloads_heads_type(),
		namespace: m.applications_workloads_heads_namespace(),
		health: m.applications_workloads_heads_health(),
		service: m.applications_workloads_heads_service(),
		pod: m.applications_workloads_heads_pod(),
		replica: m.applications_workloads_heads_replica(),
		container: m.applications_workloads_heads_container(),
		volume: m.applications_workloads_heads_volume(),
		nodeport: m.applications_workloads_heads_nodeport()
	};

	export const headers = {
		row_picker,
		name,
		type,
		namespace,
		health,
		service,
		pod,
		replica,
		container,
		volume,
		nodeport
	};
</script>

{#snippet row_picker(table: Table<Application>)}
	<Layout.Header class="justify-center">
		<Layout.HeaderController>
			<Headers.RowPicker {table} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet name(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.name}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet type(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.type}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet namespace(column: Column<Application>)}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.namespace}</Layout.HeaderViewer>
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
	</Layout.Header>
{/snippet}

{#snippet health(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.health}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet service(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.service}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet pod(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.pod}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet replica(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.replica}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet container(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.container}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet volume(column: Column<Application>)}
	<Layout.Header class="justify-end">
		<Layout.HeaderController>
			<Sorter {column} />
		</Layout.HeaderController>
		<Layout.HeaderViewer>{messages.volume}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}

{#snippet nodeport()}
	<Layout.Header class="justify-start">
		<Layout.HeaderViewer>{messages.nodeport}</Layout.HeaderViewer>
	</Layout.Header>
{/snippet}
