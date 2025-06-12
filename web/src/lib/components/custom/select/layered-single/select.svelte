<script lang="ts" module>
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
</script>

<script lang="ts">
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';

	let {
		open = $bindable(false),
		value = $bindable(),
		children,
		options,
		required,
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		value: any;
		options: OptionType[];
		required?: boolean;
	} = $props();

	const setter = (newValue: any) => {
		value = newValue;
	};
	const getter = () => {
		return value;
	};
	setContext('OptionManager', new OptionManager(options, setter, getter));
	setContext('required', required);
</script>

<DropdownMenu.Root {open} {...restProps}>
	{@render children?.()}
</DropdownMenu.Root>
