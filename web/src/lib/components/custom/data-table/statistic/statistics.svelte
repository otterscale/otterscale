<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import * as Card from '$lib/components/ui/card';
	import type { Type } from './types';
	import { setContext } from 'svelte';

	let {
		ref = $bindable(null),
		class: className,
		children,
		type,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		type: Type;
	} = $props();

	setContext('type', type);
</script>

<Card.Root
	bind:ref
	data-slot="data-table-statistics"
	class={cn('relative overflow-hidden', className)}
	{...restProps}
>
	{@render children?.()}
</Card.Root>
