<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import type { Row } from '@tanstack/table-core';
	import type { Bucket } from './types';

	export const cells = {
		_row_picker: _row_picker,
		name: name,
		owner: owner,
		usedCapacity: usedCapacity,
		capacityLimit: capacityLimit,
		objects: objects,
		objectLimit: objectLimit
	};
</script>

{#snippet _row_picker(row: Row<Bucket>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet name(row: Row<Bucket>)}
	<p>{row.original.name}</p>
{/snippet}

{#snippet owner(row: Row<Bucket>)}
	<Badge variant="outline">{row.original.owner}</Badge>
{/snippet}

{#snippet usedCapacity(row: Row<Bucket>)}
	<p class="text-right">{row.original.usedCapacity}</p>
{/snippet}

{#snippet capacityLimit(row: Row<Bucket>)}
	<p class="text-right">{row.original.capacityLimit}</p>
{/snippet}

{#snippet objects(row: Row<Bucket>)}
	<p class="text-right">{row.original.objects}</p>
{/snippet}

{#snippet objectLimit(row: Row<Bucket>)}
	<p class="text-right">{row.original.objectLimit}</p>
{/snippet}
