<script lang="ts">
	import { cn, type WithElementRef } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLParagraphElement>> = $props();
</script>

<p
	class={cn('text-muted-foreground bg-muted rounded-md p-2 text-xs font-light', className)}
	{...restProps}
>
	{@render children?.()}
</p>
