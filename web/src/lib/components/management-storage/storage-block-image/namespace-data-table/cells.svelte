<script lang="ts" module>
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
	import type { Namespace } from './types';
	import { Badge } from '$lib/components/ui/badge';

	export const cells = {
		_row_picker: _row_picker,
		totalImages: totalImages,
		pool: pool,
		namespace: namespace
	};
</script>

{#snippet _row_picker(row: Row<Namespace>)}
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
{/snippet}

{#snippet pool(row: Row<Namespace>)}
	{row.original.pool}
{/snippet}

{#snippet namespace(row: Row<Namespace>)}
	<Badge variant="outline">
		{row.original.namespace}
	</Badge>
{/snippet}

{#snippet totalImages(row: Row<Namespace>)}
	{row.original.totalImages}
{/snippet}
