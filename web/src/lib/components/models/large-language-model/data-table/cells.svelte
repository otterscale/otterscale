<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import { type LargeLangeageModel } from '../type';

	import Relation from './cell-relation.svelte';

	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { formatBigNumber } from '$lib/formatter';

	export const cells = {
		row_picker,
		name,
		gpu_cache,
		kv_cache,
		requests,
		time_to_first_token,
		relation,
	};
</script>

{#snippet row_picker(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet gpu_cache(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.gpu_cache}
	</Layout.Cell>
{/snippet}

{#snippet kv_cache(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{row.original.kv_cache}
	</Layout.Cell>
{/snippet}

{#snippet requests(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.requests)}
	</Layout.Cell>
{/snippet}

{#snippet time_to_first_token(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		{formatBigNumber(row.original.time_to_first_token)}
	</Layout.Cell>
{/snippet}

{#snippet relation(row: Row<LargeLangeageModel>)}
	<Layout.Cell class="items-end">
		<Relation model={row.original} />
	</Layout.Cell>
{/snippet}
