<script lang="ts">
	import type { HTMLAttributes } from 'svelte/elements';
	import { cn, type WithElementRef } from '$lib/utils.js';
	import Icon from '@iconify/svelte';

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
		'absolute top-0 right-0 text-8xl tracking-tight text-nowrap text-primary/5 uppercase group-hover:hidden *:[svg]:size-72',
		className
	)}
>
	<Icon {icon} />
</div>
