<script lang="ts" module>
	import { Button, buttonVariants, type ButtonProps } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	import { getContext } from 'svelte';
	import type { OptionManager } from './utils.svelte';
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
	}: ButtonProps = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Button
	bind:ref
	data-slot="select-all"
	class={cn('cursor-pointer', buttonVariants({ variant, size }), className)}
	{href}
	{type}
	{disabled}
	{...restProps}
	onclick={() => {
		optionManager.all();
	}}
>
	{@render children?.()}
</Button>
