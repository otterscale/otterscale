<script lang="ts">
	import { getContext } from 'svelte';

	import * as Command from '$lib/components/ui/command';
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';

	import { Command as CommandPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		onclick,
		option,
		...restProps
	}: CommandPrimitive.ItemProps & { option: OptionType } = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Command.Item
	bind:ref
	data-slot="select-item"
	class={cn(className)}
	onclick={(e) => {
		optionManager.handleSelect(option);
		onclick?.(e);
	}}
	{...restProps}
/>
