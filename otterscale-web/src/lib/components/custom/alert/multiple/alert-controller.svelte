<script lang="ts">
	import { getContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import Icon from '@iconify/svelte';

	import { Button } from '$lib/components/ui/button';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';

	import type { IterationManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		class: className,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLSpanElement>> = $props();

	let iterationManager: IterationManager = getContext('IterationManager');
</script>

<span
	bind:this={ref}
	data-slot="alert-controller"
	class={cn(
		className,
		'absolute top-0 right-0 flex h-full flex-col items-center justify-between gap-2 rounded-lg p-4 [&>svg]:cursor-pointer'
	)}
	{...restProps}
>
	{#if iterationManager.alerts.length > 1}
		<Icon
			icon="ph:caret-up"
			onclick={() => {
				iterationManager.previous();
			}}
		/>
		<Icon icon="ph:caret-down" onclick={() => iterationManager.next()} />
	{/if}
</span>
