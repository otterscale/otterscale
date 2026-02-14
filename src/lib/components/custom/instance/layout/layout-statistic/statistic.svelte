<script lang="ts" module>
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';
</script>

<script lang="ts">
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
	class={cn('h-60', 'ring-none border-none bg-sidebar shadow-none', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Root>
