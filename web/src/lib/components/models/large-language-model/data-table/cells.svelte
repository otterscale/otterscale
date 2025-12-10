<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import type { Row } from '@tanstack/table-core';

	import type { Model } from '$lib/api/model/v1/model_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
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
	<Layout.Cell class="items-center">
		<Cells.RowExpander {row} />
	</Layout.Cell>
{/snippet}

{#snippet row_picker(row: Row<Model>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet modelName(row: Row<Model>)}
	<Layout.Cell class="items-start">
		{row.original.id}
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Model>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<Model>)}
	<Layout.Cell class="items-start">
		{row.original.namespace}
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<Model>)}
	<Layout.Cell class="items-start">
		{row.original.status}
	</Layout.Cell>
{/snippet}

{#snippet description(row: Row<Model>)}
	<Layout.Cell class="items-start">
		<p class="max-w-[200px] truncate">{row.original.description}</p>
	</Layout.Cell>
{/snippet}

{#snippet firstDeployedAt(row: Row<Model>)}
	{#if row.original.firstDeployedAt}
		<Layout.Cell class="items-end">
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
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet lastDeployedAt(row: Row<Model>)}
	{#if row.original.lastDeployedAt}
		<Layout.Cell class="items-end">
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
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet prefill(row: Row<Model>)}
	{#if row.original.prefill}
		<Layout.Cell class="items-start">
			{row.original.prefill.vgpumemPercentage}%
			<Layout.SubCell>
				{row.original.prefill.replica}
				{m.replica()}
			</Layout.SubCell>
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet decode(row: Row<Model>)}
	{#if row.original.decode}
		<Layout.Cell class="items-start">
			{row.original.decode.vgpumemPercentage}%
			<Layout.SubCell>
				{row.original.decode.tensor}
				{m.tensor()}
			</Layout.SubCell>
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet gpuRelation(data: { row: Row<Model>; scope: string })}
	{#if data.row.original.status === 'deployed'}
		<Layout.Cell class="items-center">
			<GPURelation scope={data.scope} model={data.row.original} />
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet test(data: { row: Row<Model>; serviceUri: string; scope: string })}
	{@const readyPods = data.row.original.pods.filter((pod) => {
		const match = pod.ready.match(/^(\d+)\/(\d+)$/);
		if (!match) return false;
		return Number(match[1]) / Number(match[2]) === 1;
	})}
	{@const isReady = readyPods.length > 0}
	<Layout.Cell class="items-center">
		<Prompting
			serviceUri={data.serviceUri}
			model={data.row.original}
			scope={data.scope}
			disabled={!isReady}
		/>
	</Layout.Cell>
{/snippet}

{#snippet action(data: { row: Row<Model>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-end">
		<Actions model={data.row.original} scope={data.scope} reloadManager={data.reloadManager} />
	</Layout.Cell>
{/snippet}
