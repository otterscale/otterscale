<script lang="ts">
	import { cn } from '$lib/utils';
	import { type WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		children,
		align,
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & { align?: 'left' | 'right' } = $props();
</script>

<div
	bind:this={ref}
	class={cn('flex gap-2', align === 'right' && 'justify-end', className)}
	{...restProps}
>
	{@render children?.()}
</div>
