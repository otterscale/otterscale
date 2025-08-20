<script lang="ts">
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import type { Row } from '@tanstack/table-core';
</script>

<script lang="ts" generics="TData">
	let { row }: { row: Row<TData> } = $props();
</script>

<div class="flex items-center justify-center">
	<Checkbox
		checked={row.getIsSelected()}
		onCheckedChange={(value) => row.toggleSelected(!!value)}
		class="border-secondary-950"
		aria-label="Select row"
	/>
</div>
