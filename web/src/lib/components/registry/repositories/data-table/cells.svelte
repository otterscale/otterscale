<script lang="ts" module>
	import type { Row } from '@tanstack/table-core';

	import type { Repository } from '$lib/api/registry/v1/registry_pb';
	import { Cells } from '$lib/components/custom/data-table/core';
	import * as Layout from '$lib/components/custom/data-table/layout';
	import { ReloadManager } from '$lib/components/custom/reloader';
	import { formatCapacity } from '$lib/formatter';

	import Manifests from './cell-manifests.svelte';
	// import Actions from './cell-actions.svelte';

	export const cells = {
		row_picker,
		name,
		manifest_count,
		size_bytes,
		latest_tag
		// actions
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

{#snippet manifest_count(data: {
	row: Row<Repository>;
	scope: string;
	reloadManager: ReloadManager;
})}
	<Layout.Cell class="items-end">
		<Manifests
			repository={data.row.original}
			scope={data.scope}
			reloadManager={data.reloadManager}
		/>
	</Layout.Cell>
{/snippet}

{#snippet size_bytes(row: Row<Repository>)}
	{@const { value, unit } = formatCapacity(row.original.sizeBytes)}
	<Layout.Cell class="items-end">
		{value}
		{unit}
	</Layout.Cell>
{/snippet}

{#snippet latest_tag(row: Row<Repository>)}
	<Layout.Cell class="items-end">
		{row.original.latestTag}
	</Layout.Cell>
{/snippet}

<!-- {#snippet actions(row: Row<Repository>, scope: string, reloadManager: ReloadManager)}
	<Layout.Cell class="items-end">
		<Actions {row} {scope} {reloadManager} />
	</Layout.Cell>
{/snippet} -->
