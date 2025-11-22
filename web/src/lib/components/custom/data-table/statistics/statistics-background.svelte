<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import { cn, type WithElementRef } from '$lib/utils.js';

	let {
		icon,
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLParagraphElement>> & {
		icon: string;
	} = $props();
</script>

<div
	data-slot="data-table-statistics-background"
	class={cn(
		'absolute top-4 -right-12 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden *:[svg]:size-64',
		className
	)}
	{...restProps}
>
	<Icon {icon} />
</div>
