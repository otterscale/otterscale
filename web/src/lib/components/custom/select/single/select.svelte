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
		selectedOption,
		...restProps
	}: PopoverPrimitive.RootProps & {
		value: any;
		selectedOption?: OptionType;
	} = $props();

	setContext(
		'OptionManager',
		new OptionManager(selectedOption ?? ({} as OptionType), (o: OptionType) => {
			value = o.value;
		})
	);
</script>

<Popover.Root {open} {...restProps}>
	{@render children?.()}
</Popover.Root>
