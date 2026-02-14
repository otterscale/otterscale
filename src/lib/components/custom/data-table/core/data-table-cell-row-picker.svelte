<script lang="ts" generics="TData">
	import type { Row } from '@tanstack/table-core';

	import { Checkbox } from '$lib/components/ui/checkbox/index.js';

	let { row }: { row: Row<TData> } = $props();
</script>

<div data-slot="data-table-row-picker" class="flex items-center justify-center">
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
</div>
