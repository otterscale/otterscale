<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { Popover as PopoverPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';

	import { type ButtonVariant, buttonVariants } from '$lib/components/ui/button';
	import * as Popover from '$lib/components/ui/popover';
	import { cn } from '$lib/utils.js';

	import { type OptionManager, validate } from './utils.svelte';
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

	const required: boolean | undefined = getContext('required');
	const optionManager: OptionManager = getContext('OptionManager');

	const isInvalid = $derived(validate(required, optionManager));
</script>

<Popover.Trigger
	bind:ref
	data-slot="select-trigger"
	class={cn(
		'group cursor-pointer ring-1 data-[state=open]:ring-primary',
		buttonVariants({ variant: variant }),
		isInvalid ? 'ring-destructive' : '',
		className
	)}
	{...restProps}
>
	{#if children}
		{@render children?.()}
	{:else if optionManager.selectedOption.label}
		<div class="flex items-center gap-1 rounded-sm p-1 font-normal">
			<Icon
				icon={optionManager.selectedOption.icon ?? 'ph:empty'}
				class={cn('size-4', optionManager.selectedOption ? 'visibale' : 'hidden')}
			/>
			{optionManager.selectedOption.label}
		</div>
	{:else if isInvalid}
		<span
			class="flex items-center gap-1 text-xs group-data-[state=closed]:text-destructive group-data-[state=open]:text-primary"
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
</Popover.Trigger>
