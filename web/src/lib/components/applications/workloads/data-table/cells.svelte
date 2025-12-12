<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { Application_Type } from '$lib/api/application/v1/application_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import * as Progress from '$lib/components/custom/progress/index.js';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import Button from '$lib/components/ui/button/button.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip';

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
	const types = {
		[Application_Type.UNKNOWN]: { label: 'Unknown', icon: 'ph:question' },
		[Application_Type.DEPLOYMENT]: { label: 'Deployment', icon: 'ph:stack' },
		[Application_Type.STATEFUL_SET]: { label: 'StatefulSet', icon: 'ph:database' },
		[Application_Type.DAEMON_SET]: { label: 'DaemonSet', icon: 'ph:browsers' }
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
		{@const info = types[row.original.type] ?? types[Application_Type.UNKNOWN]}
		<Tooltip.Root>
			<Tooltip.Trigger>
				<Icon icon={info.icon} class="size-5" />
			</Tooltip.Trigger>
			<Tooltip.Content>
				{info.label}
			</Tooltip.Content>
		</Tooltip.Root>
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
			target="LTB"
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
	</Layout.Cell>
{/snippet}

{#snippet actions(data: { row: Row<Application>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-start">
		<Actions
			application={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}
