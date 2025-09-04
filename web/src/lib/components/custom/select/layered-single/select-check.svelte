<script lang="ts" module>
	import * as Command from '$lib/components/ui/command';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { OptionType } from './types';
	import type { OptionManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		option,
		parents = [],
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> & {
		option: OptionType;
		parents?: OptionType[];
	} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Command.Shortcut bind:ref data-slot="select-check" class={cn('p-0', className)} {...restProps}>
	<Icon icon="ph:check" class={optionManager.isOptionSelected(option, parents) ? 'visible' : 'invisible'} />
</Command.Shortcut>
