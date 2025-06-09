<script lang="ts" module>
	import { getContext } from 'svelte';
	import { Button, buttonVariants, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
	import type { InputManager, ValuesManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		href = undefined,
		type = 'button',
		disabled,
		variant = 'default',
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
		if (inputManager.input) {
			valuesManager.append(inputManager.input);
			inputManager.reset();
		}
	}}
	{...restProps}
>
	Add
</Button>
