<script lang="ts">
	import { getContext } from 'svelte';
	import type { OptionType } from './types';
	import type { OptionManager } from './utils.svelte';

	import * as Command from '$lib/components/ui/command';

	import { cn, type WithElementRef } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		children,
		option,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> & { option: OptionType } = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Command.Shortcut bind:ref data-slot="select-check" class={cn('p-0', className)} {...restProps}>
	<Icon icon="ph:check" class={optionManager.isOptionSelected(option) ? 'visible' : 'invisible'} />
</Command.Shortcut>
