<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import ChevronRightIcon from '@lucide/svelte/icons/chevron-right';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		inset,
		children,
		...restProps
	}: DropdownMenuPrimitive.SubTriggerProps & {
		inset?: boolean;
	} = $props();
</script>

<DropdownMenu.SubTrigger
	bind:ref
	data-slot="seelct-sub-trigger"
	data-inset={inset}
	class={cn(
		'hover:bg-accent hover:text-accent-foreground cursor-pointer rounded-sm select-none',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</DropdownMenu.SubTrigger>
