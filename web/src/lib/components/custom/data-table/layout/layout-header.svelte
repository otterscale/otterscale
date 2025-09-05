<script lang="ts" module>
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	import { cn } from '$lib/utils';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	class={cn('group flex items-center justify-between gap-1 text-xs font-medium', className)}
	{...restProps}
>
	{@render children?.()}
</div>
