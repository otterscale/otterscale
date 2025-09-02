<script lang="ts" module>
	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<Card.Action
	bind:ref
	data-slot="layout-statistic-action"
	class={cn('ml-auto flex items-center gap-1', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Action>
