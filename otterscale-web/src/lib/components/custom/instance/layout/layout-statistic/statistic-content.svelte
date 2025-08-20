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

<Card.Content
	bind:ref
	data-slot="layout-statistic-content"
	class={cn('flex h-40 items-center text-4xl', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Content>
