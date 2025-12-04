<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import Button from '$lib/components/ui/button/button.svelte';

	import type { Application } from '../types';
	import Actions from './cell-actions.svelte';

	export const cells = {
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
		nodeport,
		actions
	};
</script>

{#snippet row_picker(row: Row<Application>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

<!-- TODO: fix scope -->
{#snippet name(row: Row<Application>)}
	<Table.Cell alignClass="items-start">
		<a
			class="underline hover:no-underline"
			href={resolve('/(auth)/scope/[scope]/applications/workloads/[namespace]/[application_name]', {
				scope: page.params.scope!,
				namespace: row.original.namespace!,
				application_name: row.original.name!
			})}
		>
			{row.original.name}
		</a>
	</Table.Cell>
{/snippet}

{#snippet type(row: Row<Application>)}
	<Table.Cell alignClass="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Table.Cell>
{/snippet}

{#snippet namespace(row: Row<Application>)}
	<Table.Cell alignClass="items-start">
		{row.original.namespace}
	</Table.Cell>
{/snippet}

{#snippet health(row: Row<Application>)}
	<Table.Cell alignClass="items-end">
		<Progress.Root
			numerator={Number(row.original.healthies)}
			denominator={Number(row.original.pods.length)}
			target="LTB"
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{numerator}/{denominator}
			{/snippet}
		</Progress.Root>
	</Table.Cell>
{/snippet}

{#snippet service(row: Row<Application>)}
	<Table.Cell alignClass="items-end">
		{row.original.services.length}
	</Table.Cell>
{/snippet}

{#snippet pod(row: Row<Application>)}
	<Table.Cell alignClass="items-end">
		{row.original.pods.length}
	</Table.Cell>
{/snippet}

{#snippet replica(row: Row<Application>)}
	<Table.Cell alignClass="items-end">
		{row.original.replicas}
	</Table.Cell>
{/snippet}

{#snippet container(row: Row<Application>)}
	<Table.Cell alignClass="items-end">
		{row.original.containers.length}
	</Table.Cell>
{/snippet}

{#snippet volume(row: Row<Application>)}
	<Table.Cell alignClass="items-end">
		{row.original.persistentVolumeClaims.length}
	</Table.Cell>
{/snippet}

{#snippet nodeport(row: Row<Application>)}
	<Table.Cell alignClass="items-start">
		<div class="flex flex-wrap gap-1">
			{#each row.original.services as service (service.name)}
				{#each service.ports as port, index (index)}
					{#if port.nodePort > 0}
						<Button
							class="flex items-center"
							size="sm"
							variant="ghost"
							target="_blank"
							href={`http://${row.original.hostname}:${port.nodePort}`}
						>
							<Icon icon="ph:arrow-square-out" />
							{port.nodePort}
						</Button>
					{/if}
				{/each}
			{/each}
		</div>
	</Table.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Application>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-start">
		<Actions
			application={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Table.Cell>
{/snippet}
