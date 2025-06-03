<script lang="ts">
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	data-slot="select-actions"
	class={cn('grid grid-cols-2 gap-1 p-1 *:col-span-1', className)}
	{...restProps}
>
	{@render children?.()}
</div>
