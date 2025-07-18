<script lang="ts" module>
	import * as Popover from '$lib/components/ui/popover';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		id,
		open = $bindable(false),
		value = $bindable(),
		children,
		options = $bindable(),
		required,
		...restProps
	}: PopoverPrimitive.RootProps & {
		id?: string;
		options: Writable<OptionType[]>;
		value: any;
		required?: boolean;
	} = $props();

	setContext('options', options);
	setContext(
		'OptionManager',
		new OptionManager($options, {
			get value() {
				return value ?? '';
			},
			set value(newValue: any) {
				value = newValue;
			}
		})
	);
	setContext('id', id);
	setContext('required', required);
</script>

<Popover.Root {open} {...restProps}>
	{@render children?.()}
</Popover.Root>
