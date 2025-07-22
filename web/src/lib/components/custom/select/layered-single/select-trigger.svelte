<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import { buttonVariants } from '$lib/components/ui/button';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu';
	import Separator from '$lib/components/ui/separator/separator.svelte';
	import { cn } from '$lib/utils';
	import Icon from '@iconify/svelte';
	import { DropdownMenu as DropdownMenuPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
</script>

<script lang="ts">
	import { OptionManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		children,
		...restProps
	}: DropdownMenuPrimitive.TriggerProps & {} = $props();

	const optionManager: OptionManager = getContext('OptionManager');
	const required: boolean | undefined = getContext('required');
	const id: string | undefined = getContext('id');
	const isNotFilled = $derived(required && !optionManager.selectedAncestralOption);

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
		'cursor-pointer',
		buttonVariants({ variant: 'outline' }),
		required && isNotFilled ? 'ring-destructive ring-1' : 'ring-1'
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
			<Icon
				icon={option.icon ?? 'ph:empty'}
				class={cn(option.icon && option.icon ? 'visibale' : 'hidden')}
			/>
			{option.label}
		{/each}
	{:else if required && isNotFilled}
		<p class=" text-destructive text-xs">Required</p>
	{:else}
		Select
	{/if}
</DropdownMenu.Trigger>
