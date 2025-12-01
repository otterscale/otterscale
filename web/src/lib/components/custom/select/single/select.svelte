<script lang="ts" module>
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	import * as Popover from '$lib/components/ui/popover';

	import type { OptionType } from './types';
	import { OptionManager, validate } from './utils.svelte';
</script>

<script lang="ts">
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
		options: OptionType[];
		value?: any;
		required?: boolean;
		invalid?: boolean | null | undefined;
	} = $props();
	const optionManager = new OptionManager(options, {
		get value() {
			return value ?? '';
		},
		set value(newValue) {
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
