<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import { type LargeLangeageModel } from '../type';

	import Actions from './cell-actions.svelte';
	import Relation from './cell-relation.svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { formatBigNumber } from '$lib/formatter';
	import { dynamicPaths } from '$lib/path';

	export const cells = {
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

{#snippet row_picker(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet model(row: Row<LargeLangeageModel>)}
	<Layout.Cell
		class="items-start underline hover:cursor-pointer hover:no-underline"
		onclick={() => {
			goto(
				`${dynamicPaths.applicationsWorkloads(page.params.scope).url}/${row.original.application.namespace}/${row.original.application.name}`,
			);
		}}
	>
		{row.original.application.name}
		<Layout.SubCell>
			{row.original.application.namespace}
		</Layout.SubCell>
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet replicas(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.application.replicas}
	</Layout.Cell>
{/snippet}

{#snippet healthies(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.application.healthies}
	</Layout.Cell>
{/snippet}

{#snippet gpu_cache(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.metrics.gpu_cache}
	</Layout.Cell>
{/snippet}

{#snippet kv_cache(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.metrics.kv_cache}
	</Layout.Cell>
{/snippet}

{#snippet requests(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.metrics.requests)}
	</Layout.Cell>
{/snippet}

{#snippet time_to_first_token(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.metrics.time_to_first_token)}
	</Layout.Cell>
{/snippet}

{#snippet relation(row: Row<LargeLangeageModel>)}
	{#if row.original.application.healthies > 0}
		<Layout.Cell class="items-end">
			<Relation model={row.original} />
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet action(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		<Actions model={row.original} />
	</Layout.Cell>
{/snippet}
