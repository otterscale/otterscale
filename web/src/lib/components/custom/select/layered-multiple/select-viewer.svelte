<script lang="ts">
	import { getContext } from 'svelte';
	import type { OptionManager } from './utils.svelte';
	import Badge from '$lib/components/ui/badge/badge.svelte';
	import Icon from '@iconify/svelte';

	import { cn } from '$lib/utils.js';
	import type { HTMLAttributes } from 'svelte/elements';
	import type { WithElementRef } from 'bits-ui';
	import Separator from '$lib/components/ui/separator/separator.svelte';

	let {
		ref = $bindable(null),
		class: className,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();

	const optionManager: OptionManager = getContext('OptionManager');
</script>

{#if optionManager.selectedAncestralOptions.length}
	<div
		bind:this={ref}
		data-slot="select-viewer"
		class={cn('flex flex-wrap gap-1', className)}
		{...restProps}
	>
		{#each optionManager.selectedAncestralOptions as option}
			<Badge variant="outline" class={cn('flex items-center gap-1 rounded-sm p-1 font-normal')}>
				{#each option as part, index}
					{#if index > 0}
						<Separator class="data-[orientation=vertical]:h-3" orientation="vertical" />
					{/if}
					<Icon
						icon={part.icon ?? 'ph:empty'}
						class={cn(part.icon && part.icon ? 'visibale' : 'hidden')}
					/>
					{part.label}
				{/each}
			</Badge>
		{/each}
	</div>
{/if}
