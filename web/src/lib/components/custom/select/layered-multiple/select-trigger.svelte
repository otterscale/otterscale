<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import * as HoverCard from '$lib/components/ui/hover-card';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
</script>

<script lang="ts">
	import type { AncestralOptionType } from './types';
	import { OptionManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		children,
		...restProps
	}: DropdownMenuPrimitive.TriggerProps & {} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
	const required: boolean | undefined = getContext('required');
	const id: string | undefined = getContext('id');

	const isNotFilled = $derived(required && !optionManager.isSomeAncestralOptionsSelected);
	const formValidator: FormValidator = getContext('FormValidator');
	$effect(() => {
		if (formValidator) {
			formValidator.set(id, isNotFilled);
		}
	});
</script>

<DropdownMenu.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn(
		'w-full cursor-pointer',
		buttonVariants({ variant: 'outline' }),
		required && isNotFilled ? 'ring-destructive ring-1' : 'ring-1'
	)}
	{...restProps}
>
	{#if children}
		{@render children?.()}
	{:else if optionManager.isSomeAncestralOptionsSelected}
		<p>Select</p>
		<Separator orientation="vertical" />
		<div class="flex gap-1">
			{#if optionManager.selectedAncestralOptions.length > optionManager.visibility}
				<HoverCard.Root>
					<HoverCard.Trigger>
						{@render ShowSimplifiedOptions(optionManager.selectedAncestralOptions)}
					</HoverCard.Trigger>
					<HoverCard.Content class="flex w-fit flex-col gap-2 p-2">
						{@render ListOptions(optionManager.selectedAncestralOptions)}
					</HoverCard.Content>
				</HoverCard.Root>
			{:else}
				{@render ShowOptions(optionManager.selectedAncestralOptions)}
			{/if}
		</div>
	{:else if required && isNotFilled}
		<p class=" text-destructive text-xs">Required</p>
	{:else}
		Select
	{/if}
</DropdownMenu.Trigger>

{#snippet ShowOption(selectedOption: AncestralOptionType)}
	<Badge variant="outline" class={cn('flex items-center gap-1 rounded-sm p-1 font-normal')}>
		{#each selectedOption as part, index}
			{#if index > 0}
				<Separator orientation="vertical" />
			{/if}
			<Icon
				icon={part.icon ?? 'ph:empty'}
				class={cn(part.icon && part.icon ? 'visibale' : 'hidden')}
			/>
			{part.label}
		{/each}
	</Badge>
{/snippet}

{#snippet ShowOptions(selectedOptions: AncestralOptionType[])}
	{#each selectedOptions as option}
		{@render ShowOption(option)}
	{/each}
{/snippet}

{#snippet ShowSimplifiedOptions(selectedOptions: AncestralOptionType[])}
	<span class="flex items-center gap-1">
		<p>{selectedOptions.length}</p>
		<Icon icon="ph:check" class="size-3" />
	</span>
{/snippet}

{#snippet ListOptions(selectedOptions: AncestralOptionType[])}
	{#each selectedOptions as option}
		<span class="flex items-center gap-1 text-xs">
			{#each option as part, index}
				{#if index > 0}
					<Separator orientation="vertical" class="data-[orientation=vertical]:h-3" />
				{/if}
				<Icon
					icon={part.icon ?? 'ph:empty'}
					class={cn(part.icon && part.icon ? 'visibale' : 'hidden')}
				/>
				{part.label}
			{/each}
		</span>
	{/each}
{/snippet}
