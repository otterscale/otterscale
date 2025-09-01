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
	}: WithElementRef<HTMLAttributes<HTMLParagraphElement>> = $props();
</script>

<Card.Description bind:ref data-slot="layout-statistic-description" class={cn('', className)} {...restProps}>
	{@render children?.()}
</Card.Description>
