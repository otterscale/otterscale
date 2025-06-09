<script lang="ts">
	import { cn } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import Icon from '@iconify/svelte';

	let {
		ref = $bindable(null),
		class: className,
		isNotFilled,
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		isNotFilled: boolean | null | undefined;
	} = $props();
</script>

{#if isNotFilled}
	<div
		bind:this={ref}
		data-slot="input-required"
		class={cn('animate-in fade-in text-destructive -mb-1 flex items-center gap-1', className)}
		{...restProps}
	>
		<Icon icon="ph:asterisk" class="size-2" />
		<p class="text-xs font-light">Required</p>
	</div>
{/if}
