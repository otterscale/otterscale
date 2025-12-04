<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Repository } from '$lib/api/registry/v1/registry_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { formatCapacity } from '$lib/formatter';

	import Manifests from './cell-manifests.svelte';

	export const cells = {
		row_picker,
		name,
		manifests,
		sizeBytes,
		latestTag
	};
</script>

{#snippet row_picker(row: Row<Repository>)}
	<Layout.Cell class="items-center">
		<Cells.RowPicker {row} />
	</Layout.Cell>
{/snippet}

{#snippet name(row: Row<Repository>)}
	<Layout.Cell class="items-start">
		{row.original.name}
	</Layout.Cell>
{/snippet}

{#snippet manifests(data: { row: Row<Repository>; scope: string; reloadManager: ReloadManager })}
	<Layout.Cell class="items-end">
		<Manifests
			repository={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}

{#snippet sizeBytes(row: Row<Repository>)}
	{@const { value, unit } = formatCapacity(row.original.sizeBytes)}
	<Layout.Cell class="items-end">
		{value}
		{unit}
	</Layout.Cell>
{/snippet}

{#snippet latestTag(row: Row<Repository>)}
	<Layout.Cell class="items-start">
		{row.original.latestTag}
	</Layout.Cell>
{/snippet}
