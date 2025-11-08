<script lang="ts" module>
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';

	import type { OptionType } from '../single';

	import { OptionManager } from './utils.svelte';

	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
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
		'cursor-pointer rounded-sm select-none hover:bg-accent hover:text-accent-foreground',
		className
	)}
	{...restProps}
/>
