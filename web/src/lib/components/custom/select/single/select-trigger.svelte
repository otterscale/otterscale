<script lang="ts">
	import { getContext } from 'svelte';
	import { buttonVariants, type ButtonVariant } from '$lib/components/ui/button';
	import * as Popover from '$lib/components/ui/popover';
	import Icon from '@iconify/svelte';
	import type { OptionManager } from './utils.svelte';

	import { cn } from '$lib/utils.js';
	import { Popover as PopoverPrimitive } from 'bits-ui';

	let {
		ref = $bindable(null),
		children,
		class: className,
		variant = 'outline',
		...restProps
	}: PopoverPrimitive.TriggerProps & {
		variant?: ButtonVariant;
	} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Popover.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn('cursor-pointer', buttonVariants({ variant: variant }), className)}
	{...restProps}
>
	{#if children}
		{@render children?.()}
	{:else if optionManager.selectedOption.label}
		<div class={cn('flex items-center gap-1 rounded-sm p-1 font-normal')}>
			<Icon
				icon={optionManager.selectedOption.icon ?? 'ph:empty'}
				class={cn('size-4', optionManager.selectedOption ? 'visibale' : 'hidden')}
			/>
			{optionManager.selectedOption.label}
		</div>
	{:else}
		Select
	{/if}
</Popover.Trigger>
