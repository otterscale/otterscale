<script lang="ts" module>
	import { timestampDate } from '@bufbuild/protobuf/wkt';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	import { type LargeLanguageModel } from '../type';

	import Actions from './cell-actions.svelte';
	import Relation from './cell-relation.svelte';

	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { formatTimeAgo } from '$lib/formatter';

	export const cells = {
		row_expander,
		row_picker,
		name,
		namespace,
		chart_version,
		application_version,
		status,
		description,
		requests,
		limits,
		first_deployed_at,
		last_deployed_at,
		relation,
		pods,
		action,
	};
</script>

{#snippet row_expander(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-center">
		<button
			class="flex size-6 items-center justify-center"
			onclick={row.getToggleExpandedHandler()}
			aria-expanded={row.getIsExpanded()}
		>
			<Icon
				icon="ph:caret-left"
				class={row.getIsExpanded()
					? 'rotate-90 transition-transform duration-300'
					: '-rotate-90 transition-transform duration-300'}
			/>
		</button>
	</Layout.Cell>
{/snippet}

{#snippet row_picker(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet namespace(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		{row.original.namespace}
	</Layout.Cell>
{/snippet}

{#snippet chart_version(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.chartVersion}
	</Layout.Cell>
{/snippet}

{#snippet application_version(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.appVersion}
	</Layout.Cell>
{/snippet}

{#snippet status(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		{row.original.status}
	</Layout.Cell>
{/snippet}

{#snippet description(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-start">
		{row.original.description}
	</Layout.Cell>
{/snippet}

{#snippet requests(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.requests}
	</Layout.Cell>
{/snippet}

{#snippet limits(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.limits}
	</Layout.Cell>
{/snippet}

{#snippet first_deployed_at(row: Row<LargeLanguageModel>)}
	{#if row.original.firstDeployedAt}
		<Layout.Cell class="items-end">
			{formatTimeAgo(timestampDate(row.original.firstDeployedAt))}
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet last_deployed_at(row: Row<LargeLanguageModel>)}
	{#if row.original.lastDeployedAt}
		<Layout.Cell class="items-end">
			{formatTimeAgo(timestampDate(row.original.lastDeployedAt))}
		</Layout.Cell>
	{/if}
{/snippet}

{#snippet pods(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		{row.original.pods.length}
	</Layout.Cell>
{/snippet}

{#snippet relation(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		<Relation model={row.original} />
	</Layout.Cell>
{/snippet}

{#snippet action(row: Row<LargeLanguageModel>)}
	<Layout.Cell class="items-end">
		<Actions model={row.original} />
	</Layout.Cell>
{/snippet}
