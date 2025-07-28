<script lang="ts" module>
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		id,
		open = $bindable(false),
		value = $bindable(),
		children,
		options,
		required,
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		id?: string;
		value: any;
		options: OptionType[];
		required?: boolean;
	} = $props();

	setContext('id', id);
	setContext('required', required);
	setContext(
		'OptionManager',
		new OptionManager(options, {
			get value() {
				return value ?? '';
			},
			set value(newValue: any) {
				value = newValue;
			}
		})
	);
</script>

<DropdownMenu.Root {open} {...restProps}>
	{@render children?.()}
</DropdownMenu.Root>
