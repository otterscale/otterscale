<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { setContext } from 'svelte';
	import { OptionManager } from './utils.svelte';
	import type { AncestralOptionType } from './types';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';

	let {
		open = $bindable(false),
		value = $bindable(),
		children,
		selectedAncestralOption,
		...restProps
	}: DropdownMenuPrimitive.RootProps & {
		value: any;
		selectedAncestralOption?: AncestralOptionType;
	} = $props();

	setContext(
		'OptionManager',
		new OptionManager(
			selectedAncestralOption ?? ({} as AncestralOptionType),
			(ancestralOption: AncestralOptionType) => {
				value = ancestralOption.map((component) => component.value);
			}
		)
	);
</script>

<DropdownMenu.Root {open} {...restProps}>
	{@render children?.()}
</DropdownMenu.Root>
