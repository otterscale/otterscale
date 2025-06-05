<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { setContext } from 'svelte';
	import * as Popover from '$lib/components/ui/popover';
	import { OptionManager } from './utils.svelte';
	import type { OptionType, AncestralOptionType } from './types';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils';

	let {
		open = $bindable(false),
		value = $bindable(),
		class: className,
		children,
		options,
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		class?: string;
		value: any[];
		options: OptionType[];
	} = $props();

	const setter = (newValues: any[]) => {
		value = newValues;
	};
	const getter = () => {
		return value ?? [];
	};
	setContext('OptionManager', new OptionManager(options, setter, getter));
</script>

<DropdownMenu.Root {open} {...restProps}>
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</DropdownMenu.Root>
