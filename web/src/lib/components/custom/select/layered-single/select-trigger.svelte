<script lang="ts">
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';

	import { getContext } from 'svelte';
	import Icon from '@iconify/svelte';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { buttonVariants } from '$lib/components/ui/button';
	import { OptionManager } from './utils.svelte';

	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { cn } from '$lib/utils';

	let {
		ref = $bindable(null),
		children,
		...restProps
	}: DropdownMenuPrimitive.TriggerProps & {} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<DropdownMenu.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn('cursor-pointer', buttonVariants({ variant: 'outline' }))}
	{...restProps}
>
	{#if optionManager.selectedAncestralOption && optionManager.selectedAncestralOption.length > 0}
		{#each optionManager.selectedAncestralOption as option, index}
			{#if index > 0}
				<Separator orientation="vertical" />
			{/if}
			<Icon
				icon={option.icon ?? 'ph:empty'}
				class={cn(option.icon && option.icon ? 'visibale' : 'hidden')}
			/>
			{option.label}
		{/each}
	{:else if children}
		{@render children?.()}
	{:else}
		Select
	{/if}
</DropdownMenu.Trigger>
