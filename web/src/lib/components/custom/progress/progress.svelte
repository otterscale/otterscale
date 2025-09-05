<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import { Progress } from '$lib/components/ui/progress';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
	import { formatProgressColor } from '$lib/formatter';
	import { cn } from '$lib/utils.js';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		numerator,
		denominator,
		ratio,
		detail,
		highIsGood = true,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		numerator: number;
		denominator: number;
		ratio?: Snippet<[{ numerator: number; denominator: number }]>;
		detail?: Snippet<[{ numerator: number; denominator: number }]>;
		highIsGood?: boolean;
	} = $props();

	const progressRatio = $derived(denominator > 0 ? numerator / denominator : 0);
</script>

{#if denominator > 0}
	<div bind:this={ref} data-slot="progress-root" {...restProps}>
		<Progress value={progressRatio} max={1} class={formatProgressColor(progressRatio, false, highIsGood)} />
		{#if ratio}
			<div
				class={cn(
					'text-muted-foreground flex items-center justify-end font-light sm:min-w-[100px] md:min-w-[200px]',
					className,
				)}
			>
				{#if detail}
					<Tooltip.Provider>
						<Tooltip.Root>
							<Tooltip.Trigger>
								{@render ratio?.({ numerator, denominator })}
							</Tooltip.Trigger>
							<Tooltip.Content>
								{@render detail?.({ numerator, denominator })}
							</Tooltip.Content>
						</Tooltip.Root>
					</Tooltip.Provider>
				{:else}
					{@render ratio?.({ numerator, denominator })}
				{/if}
			</div>
		{/if}
	</div>
{:else}
	<Icon icon="ph:infinity" />
{/if}
