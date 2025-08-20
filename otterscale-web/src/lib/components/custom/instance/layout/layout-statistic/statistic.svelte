<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Card.Root
	bind:ref
	data-slot="layout-statistic"
	class={cn('h-60', 'ring-none bg-sidebar border-none shadow-none', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Root>
