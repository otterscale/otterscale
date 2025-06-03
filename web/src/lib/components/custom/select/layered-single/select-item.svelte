<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { getContext } from 'svelte';
	import { OptionManager } from './utils.svelte';

	import { cn } from '$lib/utils.js';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import type { OptionType } from '../single';

	let {
		ref = $bindable(null),
		class: className,
		inset,
		variant = 'default',
		option,
		parents = [],
		onclick,
		...restProps
	}: DropdownMenuPrimitive.ItemProps & {
		inset?: boolean;
		variant?: 'default' | 'destructive';
		option: OptionType;
		parents?: OptionType[];
	} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<DropdownMenu.Item
	bind:ref
	data-slot="select-item"
	data-inset={inset}
	data-variant={variant}
	onclick={(e) => {
		optionManager.handleSelect(option, parents);
		onclick?.(e);
	}}
	class={cn(
		'hover:bg-accent hover:text-accent-foreground cursor-pointer rounded-sm select-none',
		className
	)}
	{...restProps}
/>
