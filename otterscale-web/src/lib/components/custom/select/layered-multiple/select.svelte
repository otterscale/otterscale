<script lang="ts" module>
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { cn } from '$lib/utils';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';
	import { validate } from './utils.svelte';
</script>

<script lang="ts">
	let {
		id,
		open = $bindable(false),
		value = $bindable(),
		class: className,
		children,
		options,
		required,
		invalid = $bindable(),
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		id?: string;
		class?: string;
		value: any[];
		options: OptionType[];
		required?: boolean;
		invalid?: boolean | null | undefined;
	} = $props();
	const optionManager = new OptionManager(options, {
		get value() {
			return value ?? [];
		},
		set value(newValues: any[]) {
			value = newValues;
		}
	});
	setContext('id', id);
	setContext('required', required);
	setContext('OptionManager', optionManager);
	$effect(() => {
		invalid = validate(required, optionManager);
	});
</script>

<DropdownMenu.Root {open} {...restProps}>
	<div class={cn('flex flex-col gap-2', className)}>
		{@render children?.()}
	</div>
</DropdownMenu.Root>
