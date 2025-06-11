<script lang="ts">
	import { Progress } from '$lib/components/ui/progress';
	import { cn } from '$lib/utils.js';
	import type { WithElementRef } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className,
		numerator,
		denominator,
		ratio,
		detail,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		numerator: number;
		denominator: number;
		ratio?: Snippet<[{ numerator: number; denominator: number }]>;
		detail?: Snippet<[{ numerator: number; denominator: number }]>;
	} = $props();
</script>

<div bind:this={ref} data-slot="progress-root" {...restProps}>
	<Progress value={numerator / denominator} max={1} />
	{#if detail || ratio}
		<div
			class={cn(
				'text-muted-foreground flex items-center justify-between gap-4 font-light',
				className
			)}
		>
			<span>
				{@render detail?.({ numerator, denominator })}
			</span>
			<span>
				{@render ratio?.({ numerator, denominator })}
			</span>
		</div>
	{/if}
</div>
