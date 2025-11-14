<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import { getContext } from 'svelte';

	import { Button, type ButtonProps, buttonVariants } from '$lib/components/ui/button/index.js';
	import { cn } from '$lib/utils.js';

	import type { ValuesManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		values,
		class: className,
		href = undefined,
		type = 'button',
		disabled,
		variant = 'outline',
		size = 'sm',
		...restProps
	}: ButtonProps & { values: any[] } = $props();

	const valuesManager: ValuesManager = getContext('ValuesManager');
</script>

<Button
	bind:ref
	data-slot="input-trigger"
	class={cn('w-fit cursor-pointer shadow', buttonVariants({ variant, size }), className)}
	{href}
	{type}
	{disabled}
	onclick={() => {
		valuesManager.values = values ?? [];
	}}
	{...restProps}
>
	<Icon icon="ph:arrow-square-in" class="size-5 text-primary" />
</Button>
