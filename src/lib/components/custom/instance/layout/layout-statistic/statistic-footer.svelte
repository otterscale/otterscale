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

<Card.Footer
	bind:ref
	data-slot="layout-statistic-footer"
	class={cn('h-10 items-end truncate text-xs font-light text-muted-foreground', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Footer>
