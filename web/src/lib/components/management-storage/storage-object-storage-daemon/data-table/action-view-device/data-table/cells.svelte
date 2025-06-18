<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { Device } from './types';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		_row_picker: _row_picker,
		id: id,
		name: name,
		state: state,
		lifeExpectancy: lifeExpectancy,
		daemons: daemons
	};
</script>

{#snippet _row_picker(row: Row<Device>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet id(row: Row<Device>)}
	{row.original.id}
{/snippet}

{#snippet name(row: Row<Device>)}
	{row.original.name}
{/snippet}

{#snippet state(row: Row<Device>)}
	<Badge variant="outline">
		{row.original.state}
	</Badge>
{/snippet}

{#snippet lifeExpectancy(row: Row<Device>)}
	<Badge variant="outline">
		{row.original.lifeExpectancy}
	</Badge>
{/snippet}

{#snippet daemons(row: Row<Device>)}
	<span class="flex items-center gap-1">
		{#each row.original.daemons as daemon}
			<Badge variant="outline">
				{daemon}
			</Badge>
		{/each}
	</span>
{/snippet}
