<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { setContext } from 'svelte';
	import { OptionManager } from './utils.svelte';
	import type { OptionType } from './types';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	let {
		open = $bindable(false),
		value = $bindable(),
		children,
		options,
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		value: any;
		options: OptionType[];
	} = $props();

	const setter = (newValue: any) => {
		value = newValue;
	};
	const getter = () => {
		return value;
	};

	setContext('OptionManager', new OptionManager(options, setter, getter));
</script>

<DropdownMenu.Root {open} {...restProps}>
	{@render children?.()}
</DropdownMenu.Root>
