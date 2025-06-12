<script lang="ts" module>
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
</script>

<script lang="ts">
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';

	let {
		open = $bindable(false),
		value = $bindable(),
		class: className,
		children,
		options,
		required,
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		class?: string;
		value: any[];
		options: OptionType[];
		required?: boolean;
	} = $props();

	const setter = (newValues: any[]) => {
		value = newValues;
	};
	const getter = () => {
		return value ?? [];
	};

	setContext('OptionManager', new OptionManager(options, setter, getter));
	setContext('required', required);
</script>

<DropdownMenu.Root {open} {...restProps}>
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</DropdownMenu.Root>
