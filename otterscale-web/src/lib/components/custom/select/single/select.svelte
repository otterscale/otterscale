<script lang="ts">
	import * as Popover from '$lib/components/ui/popover';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { OptionType } from './types';
	import { OptionManager, validate } from './utils.svelte';

	let {
		id,
		open = $bindable(false),
		value = $bindable(),
		children,
		options = $bindable(),
		required,
		invalid = $bindable(),
		...restProps
	}: PopoverPrimitive.RootProps & {
		id?: string;
		options: Writable<OptionType[]>;
		value?: any;
		required?: boolean;
		invalid?: boolean | null | undefined;
	} = $props();
	const optionManager = new OptionManager($options, {
		get value() {
			return value ?? '';
		},
		set value(newValue: any) {
			value = newValue;
		}
	});
	setContext('id', id);
	setContext('required', required);
	setContext('options', options);
	setContext('OptionManager', optionManager);
	$effect(() => {
		invalid = validate(required, optionManager);
	});
</script>

<Popover.Root {open} {...restProps}>
	{@render children?.()}
</Popover.Root>
