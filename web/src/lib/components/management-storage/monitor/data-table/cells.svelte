<script lang="ts" module>
	import { goto } from '$app/navigation';
	import type { MON } from '$gen/api/storage/v1/storage_pb';
	import TableRowPicker from '$lib/components/custom/data-table/data-table-row-pickers/cell.svelte';
	import { Badge } from '$lib/components/ui/badge';
	import Icon from '@iconify/svelte';
	import type { Row } from '@tanstack/table-core';

	export const cells = {
		_row_picker: _row_picker,
		leader: leader,
		name: name,
		rank: rank,
		publicAddress: publicAddress,
		machine: machine
	};
</script>

{#snippet _row_picker(row: Row<MON>)}
	<TableRowPicker {row} />
{/snippet}

{#snippet leader(row: Row<MON>)}
	{#if row.original.leader}
		<Icon icon="ph:circle" class="text-primary" />
	{/if}
{/snippet}

{#snippet name(row: Row<MON>)}
	{row.original.name}
{/snippet}

{#snippet rank(row: Row<MON>)}
	<div class="text-end">{row.original.rank}</div>
{/snippet}

{#snippet publicAddress(row: Row<MON>)}
	<Badge variant="outline">
		{row.original.publicAddress}
	</Badge>
{/snippet}

{#snippet machine(row: Row<MON>)}
	<div class="flex items-center gap-1">
		<Badge variant="outline">
			{row.original.machine?.hostname}
		</Badge>
		<Icon
			icon="ph:arrow-square-out"
			class="hover:cursor-pointer"
			onclick={() => {
				goto(`/management/machine/${row.original.machine?.id}`);
			}}
		/>
	</div>
{/snippet}
