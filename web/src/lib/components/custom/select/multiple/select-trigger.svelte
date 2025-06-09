<script lang="ts">
	import { buttonVariants, type ButtonVariant } from '$lib/components/ui/button';
	import * as Popover from '$lib/components/ui/popover';
	import Icon from '@iconify/svelte';
	import type { OptionType } from './types';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { getContext } from 'svelte';
	import type { OptionManager } from './utils.svelte';

	import { cn } from '$lib/utils.js';
	import { Popover as PopoverPrimitive } from 'bits-ui';

	let {
		ref = $bindable(null),
		children,
		variant = 'outline',
		class: className,
		...restProps
	}: PopoverPrimitive.TriggerProps & {
		label?: string;
		variant?: ButtonVariant;
	} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Popover.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn('w-full cursor-pointer', buttonVariants({ variant: variant }), className)}
	{...restProps}
>
	{#if children}
		{@render children?.()}
	{:else}
		Select
	{/if}
	{#if optionManager.isSomeOptionsSelected}
		<Separator orientation="vertical" />
		<div class="flex gap-1">
			{#if optionManager.selectedOptions.length > optionManager.visibility}
				<HoverCard.Root>
					<HoverCard.Trigger>
						{@render ShowSimplifiedOptions(optionManager.selectedOptions)}
					</HoverCard.Trigger>
					<HoverCard.Content class="flex w-fit flex-col gap-2 p-2">
						{@render ListOptions(optionManager.selectedOptions)}
					</HoverCard.Content>
				</HoverCard.Root>
			{:else}
				{@render ShowOptions(optionManager.selectedOptions)}
			{/if}
		</div>
	{/if}
</Popover.Trigger>

{#snippet ShowOption(selectedOption: OptionType)}
	<Badge variant="outline" class={cn('flex items-center gap-1 rounded-sm p-1 font-normal')}>
		<Icon
			icon={selectedOption.icon ? selectedOption.icon : 'ph:empty'}
			class={cn('size-3', selectedOption.icon ? 'visibale' : 'hidden')}
		/>
		{selectedOption.label}
	</Badge>
{/snippet}

{#snippet ShowOptions(selectedOptions: OptionType[])}
	{#each selectedOptions as option}
		{@render ShowOption(option)}
	{/each}
{/snippet}

{#snippet ShowSimplifiedOptions(selectedOptions: OptionType[])}
	<span class="flex items-center gap-1">
		<p>{selectedOptions.length}</p>
		<Icon icon="ph:check" class="size-3" />
	</span>
{/snippet}

{#snippet ListOptions(selectedOptions: OptionType[])}
	{#each selectedOptions as option}
		<span class="flex items-center gap-1">
			<Icon
				icon={option.icon ? option.icon : 'ph:empty'}
				class={cn('size-3', option.icon ? 'visibale' : 'invisible')}
			/>
			<p class="text-xs select-none">{option.label}</p>
		</span>
	{/each}
{/snippet}
