<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Repository } from '$lib/api/registry/v1/registry_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import * as Table from '$lib/components/custom/table/index.js';
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
	<Table.Cell alignClass="items-center">
		<Cells.RowPicker {row} />
	</Table.Cell>
{/snippet}

{#snippet name(row: Row<Repository>)}
	<Table.Cell alignClass="items-start">
		{row.original.name}
	</Table.Cell>
{/snippet}

{#snippet manifests(data: { row: Row<Repository>; scope: string; reloadManager: ReloadManager })}
	<Table.Cell alignClass="items-end">
		<Manifests
			repository={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Table.Cell>
{/snippet}

{#snippet sizeBytes(row: Row<Repository>)}
	{@const { value, unit } = formatCapacity(row.original.sizeBytes)}
	<Table.Cell alignClass="items-end">
		{value}
		{unit}
	</Table.Cell>
{/snippet}

{#snippet latestTag(row: Row<Repository>)}
	<Table.Cell alignClass="items-start">
		{row.original.latestTag}
	</Table.Cell>
{/snippet}
