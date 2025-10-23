<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import { type LargeLanguageModel } from '../type';

	import Actions from './cell-actions.svelte';
	import Relation from './cell-relation.svelte';

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

{#snippet row_picker(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet model(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		<a
			class="m-0 p-0 underline hover:no-underline"
			href={`${dynamicPaths.applicationsWorkloads(page.params.scope).url}/${row.original.application.namespace}/${row.original.application.name}`}
		>
			{row.original.application.name}
		</a>
		<Layout.SubCell>
			{row.original.application.namespace}
		</Layout.SubCell>
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet replicas(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.application.replicas}
	</Layout.Cell>
{/snippet}

{#snippet healthies(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.application.healthies}
	</Layout.Cell>
{/snippet}

{#snippet gpu_cache(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.metrics.gpu_cache}
	</Layout.Cell>
{/snippet}

{#snippet kv_cache(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.metrics.kv_cache}
	</Layout.Cell>
{/snippet}

{#snippet requests(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.metrics.requests)}
	</Layout.Cell>
{/snippet}

{#snippet time_to_first_token(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.metrics.time_to_first_token)}
	</Layout.Cell>
{/snippet}

{#snippet relation(row: Row<LargeLanguageModel>)}
	{#if row.original.application.healthies > 0}
		<Layout.Cell class="items-end">
			<Relation model={row.original} />
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet action(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		<Actions model={row.original} />
	</Layout.Cell>
{/snippet}
