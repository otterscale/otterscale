<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Card.Content
	bind:ref
	data-slot="data-table-statistics-content"
	class={cn('lg:max-[1100px]:flex-col lg:max-[1100px]:items-start', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Content>
