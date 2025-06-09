<script lang="ts">
	import { setContext } from 'svelte';
	import * as Popover from '$lib/components/ui/popover';
	import { OptionManager } from './utils.svelte';
	import type { OptionType } from './types';

	import { Popover as PopoverPrimitive } from 'bits-ui';

	let {
		open = $bindable(false),
		value = $bindable(),
		children,
		options,
		required,
		...restProps
	}: PopoverPrimitive.RootProps & {
		options: OptionType[];
		value: any;
		required?: boolean;
	} = $props();

	const setter = (o: OptionType) => {
		value = o.value;
	};
	const getter = () => {
		return value ?? '';
	};
	setContext('OptionManager', new OptionManager(options, setter, getter));
	setContext('required', required);
</script>

<Popover.Root {open} {...restProps}>
	{@render children?.()}
</Popover.Root>
