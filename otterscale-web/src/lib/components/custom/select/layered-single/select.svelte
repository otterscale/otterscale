<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { setContext } from 'svelte';
	import type { OptionType } from './types';
	import { OptionManager, validate } from './utils.svelte';

	let {
		id,
		open = $bindable(false),
		value = $bindable(),
		children,
		options,
		required,
		invalid = $bindable(),
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		id?: string;
		value: any;
		options: OptionType[];
		required?: boolean;
		invalid?: boolean | null | undefined;
	} = $props();
	const optionManager = new OptionManager(options, {
		get value() {
			return value ?? '';
		},
		set value(newValue: any) {
			value = newValue;
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
	{@render children?.()}
</DropdownMenu.Root>
