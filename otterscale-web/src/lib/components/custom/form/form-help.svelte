<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
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
