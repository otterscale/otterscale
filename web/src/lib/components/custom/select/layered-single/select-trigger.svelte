<script lang="ts" module>
	import { buttonVariants } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
	import { OptionManager, validate } from './utils.svelte';
</script>

<script lang="ts">
	let { ref = $bindable(null), children, ...restProps }: DropdownMenuPrimitive.TriggerProps & {} = $props();

	const required: boolean | undefined = getContext('required');
	const optionManager: OptionManager = getContext('OptionManager');

	const isInvalid = $derived(validate(required, optionManager));
</script>

<DropdownMenu.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn(
		'data-[state=open]:ring-primary group cursor-pointer',
		buttonVariants({ variant: 'outline' }),
		isInvalid ? 'ring-destructive ring-1' : 'ring-1',
	)}
	{...restProps}
>
	{#if children}
		{@render children?.()}
	{:else if optionManager.selectedAncestralOption && optionManager.selectedAncestralOption.length > 0}
		{#each optionManager.selectedAncestralOption as option, index}
			{#if index > 0}
				<Separator orientation="vertical" />
			{/if}
			<Icon icon={option.icon ?? 'ph:empty'} class={cn(option.icon && option.icon ? 'visibale' : 'hidden')} />
			{option.label}
		{/each}
	{:else if isInvalid}
		<span
			class="group-data-[state=open]:text-primary group-data-[state=closed]:text-destructive flex items-center gap-1 text-xs"
		>
			<Icon icon="ph:list" />
			<p class="group-data-[state=closed]:hidden">Select</p>
			<p class="group-data-[state=open]:hidden">Required</p>
		</span>
	{:else}
		<span class="flex items-center gap-1 text-xs">
			<Icon icon="ph:list" />
			Select
		</span>
	{/if}
</DropdownMenu.Trigger>
