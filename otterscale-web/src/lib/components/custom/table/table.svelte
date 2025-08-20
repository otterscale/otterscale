<script lang="ts">
	import * as Table from '$lib/components/ui/table';
	import type { HTMLTableAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLTableAttributes> = $props();
</script>

<Table.Root bind:ref data-slot="table" class={cn(className)} {...restProps}>
	{@render children?.()}
</Table.Root>
