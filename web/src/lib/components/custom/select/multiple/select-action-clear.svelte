<script lang="ts" module>
	import { getContext } from 'svelte';

	import { Button, type ButtonProps,buttonVariants } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';

	import type { OptionManager } from './utils.svelte';
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
	}: ButtonProps = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

<Button
	bind:ref
	data-slot="select-clear"
	class={cn('cursor-pointer capitalize', buttonVariants({ variant, size }), className)}
	{href}
	{type}
	{disabled}
	{...restProps}
	onclick={() => {
		optionManager.clear();
	}}
>
	{@render children?.()}
</Button>
