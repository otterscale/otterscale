<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress/index.js';
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
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

<!-- TODO: fix scope -->
{#snippet name(row: Row<Application>)}
	<Layout.Cell class="items-start">
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
	</Layout.Cell>
{/snippet}

{#snippet type(row: Row<Application>)}
	<Layout.Cell class="items-start">
		<Badge variant="outline">
			{row.original.type}
		</Badge>
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<Application>)}
	<Layout.Cell class="items-start">
		{row.original.namespace}
	</Layout.Cell>
{/snippet}

{#snippet health(row: Row<Application>)}
	<Layout.Cell class="items-end">
		<Progress.Root
			numerator={Number(row.original.healthies)}
			denominator={Number(row.original.pods.length)}
		>
			{#snippet ratio({ numerator, denominator })}
				{Progress.formatRatio(numerator, denominator)}
			{/snippet}
			{#snippet detail({ numerator, denominator })}
				{numerator}/{denominator}
			{/snippet}
		</Progress.Root>
	</Layout.Cell>
{/snippet}

{#snippet service(row: Row<Application>)}
	<Layout.Cell class="items-end">
		{row.original.services.length}
	</Layout.Cell>
{/snippet}

{#snippet pod(row: Row<Application>)}
	<Layout.Cell class="items-end">
		{row.original.pods.length}
	</Layout.Cell>
{/snippet}

{#snippet replica(row: Row<Application>)}
	<Layout.Cell class="items-end">
		{row.original.replicas}
	</Layout.Cell>
{/snippet}

{#snippet container(row: Row<Application>)}
	<Layout.Cell class="items-end">
		{row.original.containers.length}
	</Layout.Cell>
{/snippet}

{#snippet volume(row: Row<Application>)}
	<Layout.Cell class="items-end">
		{row.original.persistentVolumeClaims.length}
	</Layout.Cell>
{/snippet}

{#snippet nodeport(row: Row<Application>)}
	<Layout.Cell class="items-start">
		{#each row.original.services as service}
			{#each service.ports as port}
				{#if port.nodePort > 0}
					<span class="flex items-center">
						<Badge variant="outline">{port.targetPort}</Badge>
						<Button
							variant="ghost"
							target="_blank"
							href={`http://${row.original.publicAddress}:${port.nodePort}`}
						>
							<Icon icon="ph:arrow-square-out" />
						</Button>
					</span>
				{/if}
			{/each}
		{/each}
	</Layout.Cell>
{/snippet}

{#snippet actions(row: Row<Application>)}
	<Layout.Cell class="items-start">
		<Actions application={row.original} />
	</Layout.Cell>
{/snippet}
