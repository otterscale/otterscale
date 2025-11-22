<script lang="ts" module>
	function getIcon(type: Type) {
		switch (type) {
			case 'count':
				return 'ph:chart-bar-bold';
			case 'ratio':
				return 'ph:chart-pie-bold';
			default:
				return 'ph:chart-scatter-bold';
		}
	}
</script>

<script lang="ts">
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import * as Card from '$lib/components/ui/card';
	import { cn, type WithElementRef } from '$lib/utils.js';

	import type { Type } from './types';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();

	const type: Type = getContext('type');
</script>

<Card.Title
	bind:ref
	data-slot="data-table-statistics-title"
	class={cn('flex items-center gap-2 font-bold', className)}
	{...restProps}
>
	<div class="flex shrink-0 items-center justify-center rounded-md bg-primary/10 p-2 text-primary">
		<Icon data-slot="data-table-statistics-title-icon" icon={getIcon(type)} class="size-5" />
	</div>
	{@render children?.()}
</Card.Title>
