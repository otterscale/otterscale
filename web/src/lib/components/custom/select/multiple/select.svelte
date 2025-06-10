<script lang="ts">
	import { setContext } from 'svelte';
	import * as Popover from '$lib/components/ui/popover';
	import { OptionManager } from './utils.svelte';
	import type { OptionType } from './types';
	import type { Writable } from 'svelte/store';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils';

	let {
		open = $bindable(false),
		class: className,
		value = $bindable(),
		children,
		options = $bindable(),
		selectedOptions,
		...restProps
	}: PopoverPrimitive.RootProps & {
		class?: string;
		value: any[];
		options: Writable<OptionType[]>;
		selectedOptions?: OptionType[];
	} = $props();

	const setter = (newValues: any[]) => {
		value = newValues;
	};
	const getter = () => {
		return Array.isArray(value) ? value : value ? [value] : [];
	};
	setContext('options', options);
	setContext('OptionManager', new OptionManager($options, setter, getter));
</script>

<Popover.Root {open} {...restProps}>
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</Popover.Root>
