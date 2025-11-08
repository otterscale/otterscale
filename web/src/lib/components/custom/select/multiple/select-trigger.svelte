<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';

	import type { OptionType } from './types';
	import { validate, type OptionManager } from './utils.svelte';

	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { buttonVariants, type ButtonVariant } from '$lib/components/ui/button';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import * as Popover from '$lib/components/ui/popover';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
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

	const required: boolean | undefined = getContext('required');
	const optionManager: OptionManager = getContext('OptionManager');

	const isInvalid = $derived(validate(required, optionManager));
</script>

<Popover.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn(
		'group w-full cursor-pointer data-[state=open]:ring-primary',
		buttonVariants({ variant: variant }),
		isInvalid ? 'ring-1 ring-destructive' : 'ring-1',
		className
	)}
	{...restProps}
>
	{#if children}
		{@render children?.()}
	{:else if isInvalid}
		<span
			class="flex items-center gap-1 text-xs group-data-[state=closed]:text-destructive group-data-[state=open]:text-primary"
		>
			<Icon icon="ph:list" />
			<p class="group-data-[state=closed]:hidden">Select</p>
			<p class="group-data-[state=open]:hidden">Required</p>
		</span>
	{:else if optionManager.isSomeOptionsSelected}
		<p>Select</p>

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
	{:else}
		Select
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
