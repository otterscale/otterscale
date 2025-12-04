<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
	import Prompting from '$lib/components/prompting/index.svelte';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatTimeAgo } from '$lib/formatter';
	import { m } from '$lib/paraglide/messages';

	import Actions from './cell-actions.svelte';
	import GPURelation from './cell-gpu-relation.svelte';

	export const cells = {
		row_expander,
		row_picker,
		name,
		modelName,
		namespace,
		status,
		description,
		prefill,
		decode,
		firstDeployedAt,
		lastDeployedAt,
		gpuRelation,
		test,
		action
	};
</script>

{#snippet row_expander(row: Row<Model>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowExpander {row} />
	</Table.Cell>
{/snippet}

{#snippet row_picker(row: Row<Model>)}
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet modelName(row: Row<Model>)}
	<Table.Cell alignClass="items-start">
		{row.original.id}
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Model>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet namespace(row: Row<Model>)}
	<Table.Cell alignClass="items-start">
		{row.original.namespace}
	</Table.Cell>
{/snippet}

{#snippet status(row: Row<Model>)}
	<Table.Cell alignClass="items-start">
		{row.original.status}
	</Table.Cell>
{/snippet}

{#snippet description(row: Row<Model>)}
	<Table.Cell alignClass="items-start">
		<p class="max-w-[200px] truncate">{row.original.description}</p>
	</Table.Cell>
{/snippet}

{#snippet firstDeployedAt(row: Row<Model>)}
	{#if row.original.firstDeployedAt}
		<Table.Cell alignClass="items-end">
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.firstDeployedAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.firstDeployedAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet lastDeployedAt(row: Row<Model>)}
	{#if row.original.lastDeployedAt}
		<Table.Cell alignClass="items-end">
			<Tooltip.Provider>
				<Tooltip.Root>
					<Tooltip.Trigger>
						{formatTimeAgo(timestampDate(row.original.lastDeployedAt))}
					</Tooltip.Trigger>
					<Tooltip.Content>
						{timestampDate(row.original.lastDeployedAt)}
					</Tooltip.Content>
				</Tooltip.Root>
			</Tooltip.Provider>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet prefill(row: Row<Model>)}
	{#if row.original.prefill}
		<Table.Cell alignClass="items-start">
			{row.original.prefill.vgpumemPercentage}%
			<Table.SubCell>
				{row.original.prefill.replica}
				{m.replica()}
			</Table.SubCell>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet decode(row: Row<Model>)}
	{#if row.original.decode}
		<Table.Cell alignClass="items-start">
			{row.original.decode.vgpumemPercentage}%
			<Table.SubCell>
				{row.original.decode.tensor}
				{m.tensor()}
			</Table.SubCell>
		</Table.Cell>
	{/if}
{/snippet}

{#snippet gpuRelation(data: { row: Row<Model>; scope: string })}
	{#if data.row.original.status === 'deployed'}
		<Table.Cell alignClass="items-center">
			<GPURelation scope={data.scope} model={data.row.original} />
		</Table.Cell>
	{/if}
{/snippet}

{#snippet test(data: { row: Row<Model>; serviceUri: string })}
	{#if data.row.original.status === 'deployed'}
		<Table.Cell alignClass="items-center">
			<Prompting serviceUri={data.serviceUri} model={data.row.original} />
		</Table.Cell>
	{/if}
{/snippet}

{#snippet action(data: { row: Row<Model>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-end">
		<Actions model={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Table.Cell>
{/snippet}
