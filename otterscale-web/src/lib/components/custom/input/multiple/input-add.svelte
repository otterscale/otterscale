<script lang="ts">
	import { Button, buttonVariants, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';
	import type { InputManager, ValuesManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		class: className,
		href = undefined,
		type = 'button',
		disabled,
		variant = 'outline',
		size = 'sm',
		children,
		...restProps
	}: ButtonProps & {} = $props();

	const inputManager: InputManager = getContext('InputManager');
	const valuesManager: ValuesManager = getContext('ValuesManager');
</script>

<Button
	bind:ref
	data-slot="input-trigger"
	class={cn('w-fit cursor-pointer shadow', buttonVariants({ variant, size }), className)}
	{href}
	{type}
	{disabled}
	onclick={(e) => {
		if (inputManager.input !== '') {
			valuesManager.append(inputManager.input);
			inputManager.reset();
		}
	}}
	{...restProps}
>
	<Icon icon="ph:plus-circle" class="text-primary size-5" />
</Button>
