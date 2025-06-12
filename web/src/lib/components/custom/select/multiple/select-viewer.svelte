<script lang="ts" module>
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import { getContext } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';
</script>

<script lang="ts">
	import type { OptionManager } from './utils.svelte';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

{#if optionManager.selectedOptions.length}
	<div
		bind:this={ref}
		data-slot="select-viewer"
		class={cn('flex flex-wrap items-center gap-1', className)}
		{...restProps}
	>
		{#each optionManager.selectedOptions as option}
			<Badge variant="outline" class={cn('flex items-center gap-1 rounded-sm p-1 font-normal')}>
				<Icon
					icon={option.icon ? option.icon : 'ph:empty'}
					class={cn('size-3', option.icon ? 'visibale' : 'hidden')}
				/>
				{option.label}
			</Badge>
		{/each}
	</div>
{/if}
