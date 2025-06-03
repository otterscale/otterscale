<script lang="ts">
	import * as Command from '$lib/components/ui/command';

	import { Command as CommandPrimitive } from 'bits-ui';
	import SearchIcon from '@lucide/svelte/icons/search';
	import { cn } from '$lib/utils.js';

	let {
		ref = $bindable(null),
		class: className,
		value = $bindable(''),
		...restProps
	}: CommandPrimitive.InputProps = $props();
</script>

<Command.Input data-slot="select-input" class={cn(className)} bind:ref {...restProps} bind:value />
