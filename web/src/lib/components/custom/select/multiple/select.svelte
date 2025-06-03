<script lang="ts">
	import { setContext } from 'svelte';
	import * as Popover from '$lib/components/ui/popover';
	import { OptionManager } from './utils.svelte';
	import type { OptionType } from './types';

	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils';

	let {
		open = $bindable(false),
		class: className,
		value = $bindable(),
		children,
		options,
		selectedOptions,
		...restProps
	}: PopoverPrimitive.RootProps & {
		class?: string;
		value: any[];
		options: OptionType[];
		selectedOptions?: OptionType[];
	} = $props();

	setContext(
		'OptionManager',
		new OptionManager(options, selectedOptions ?? ([] as OptionType[]), (options: OptionType[]) => {
			value = options.map((option) => option.value);
		})
	);
</script>

<Popover.Root {open} {...restProps}>
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</Popover.Root>
