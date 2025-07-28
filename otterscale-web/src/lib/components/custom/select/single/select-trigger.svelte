<script lang="ts" module>
	import { FormValidator } from '$lib/components/custom/form';
	import { buttonVariants, type ButtonVariant } from '$lib/components/ui/button';
	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { OptionManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		children,
		class: className,
		variant = 'outline',
		...restProps
	}: PopoverPrimitive.TriggerProps & {
		variant?: ButtonVariant;
	} = $props();

	const id: string | undefined = getContext('id');
	const required: boolean | undefined = getContext('required');
	const optionManager: OptionManager = getContext('OptionManager');

	const isInvalid = $derived(required && !optionManager.selectedOption.value);
	const formValidator: FormValidator = getContext('FormValidator');
	$effect(() => {
		formValidator.set(id, isInvalid);
	});
</script>

<Popover.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn(
		'cursor-pointer ring-1',
		buttonVariants({ variant: variant }),
		required && isInvalid ? 'ring-destructive' : '',
		className
	)}
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
	{:else if required && isInvalid}
		<p class="text-destructive text-xs">Required</p>
	{:else}
		Select
	{/if}
</Popover.Trigger>
