<script lang="ts">
	import { setContext } from 'svelte';
	import * as Popover from '$lib/components/ui/popover';
	import { OptionManager } from './utils.svelte';
	import type { OptionType } from './types';
	import type { Writable } from 'svelte/store';
	import { Popover as PopoverPrimitive } from 'bits-ui';

	let {
		open = $bindable(false),
		value = $bindable(),
		children,
		options = $bindable(),
		required,
		...restProps
	}: PopoverPrimitive.RootProps & {
		options: Writable<OptionType[]>;
		value: any;
		required?: boolean;
	} = $props();

	const setter = (newValue: any) => {
		value = newValue;
	};
	const getter = () => {
		return value ?? '';
	};
	setContext('options', options);
	setContext('OptionManager', new OptionManager($options, setter, getter));
	setContext('required', required);
</script>

<Popover.Root {open} {...restProps}>
	{@render children?.()}
</Popover.Root>
