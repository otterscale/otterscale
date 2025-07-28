<script lang="ts" module>
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import { Command as CommandPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { OptionType } from './types';
	import { OptionManager } from './utils.svelte';
</script>

<script lang="ts">
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
