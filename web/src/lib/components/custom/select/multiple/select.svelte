<script lang="ts">
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { Writable } from 'svelte/store';

	import type { OptionType } from './types';
	import { OptionManager, validate } from './utils.svelte';

	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils';

	let {
		id,
		open = $bindable(false),
		class: className,
		value = $bindable(),
		children,
		options = $bindable(),
		required,
		invalid = $bindable(),
		...restProps
	}: PopoverPrimitive.RootProps & {
		id?: string;
		class?: string;
		value: any[];
		options: Writable<OptionType[]>;
		required?: boolean;
		invalid?: boolean | null | undefined;
	} = $props();
	const optionManager = new OptionManager($options, {
		get value() {
			return Array.isArray(value) ? value : value ? [value] : [];
		},
		set value(newValues: any[]) {
			value = newValues;
		},
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
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</Popover.Root>
