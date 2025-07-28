<script lang="ts">
	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { Writable } from 'svelte/store';
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';

	let {
		id,
		open = $bindable(false),
		class: className,
		value = $bindable(),
		children,
		options = $bindable(),
		required,
		selectedOptions,
		...restProps
	}: PopoverPrimitive.RootProps & {
		id?: string;
		class?: string;
		value: any[];
		options: Writable<OptionType[]>;
		selectedOptions?: OptionType[];
		required?: boolean;
	} = $props();

	setContext('id', id);
	setContext('required', required);
	setContext('options', options);
	setContext(
		'OptionManager',
		new OptionManager($options, {
			get value() {
				return Array.isArray(value) ? value : value ? [value] : [];
			},
			set value(newValues: any[]) {
				value = newValues;
			}
		})
	);
</script>

<Popover.Root {open} {...restProps}>
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</Popover.Root>
