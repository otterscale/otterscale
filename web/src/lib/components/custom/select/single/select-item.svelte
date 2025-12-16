<script lang="ts" module>
	import { Command as CommandPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';

	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';

	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		onclick,
		onSelect,
		option,
		...restProps
	}: CommandPrimitive.ItemProps & { option: OptionType } = $props();

	const optionManager: OptionManager = getContext('OptionManager');
	const selectContext: { close: () => void } = getContext('selectContext');
</script>

<Command.Item
	bind:ref
	data-slot="select-item"
	class={cn(className)}
	onclick={(e) => {
		optionManager.handleSelect(option);
		selectContext.close();
		onclick?.(e);
	}}
	onSelect={() => {
		optionManager.handleSelect(option);
		selectContext.close();
		onSelect?.();
	}}
	{...restProps}
/>
