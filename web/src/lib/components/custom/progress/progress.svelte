<script lang="ts" module>
	import Icon from '@iconify/svelte';
	import type { WithElementRef } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	import {
		formatProgressColor,
		type ProgressTargetType
	} from '$lib/components/custom/progress/utils.svelte';
	import { Progress } from '$lib/components/ui/progress';
	import * as Tooltip from '$lib/components/ui/tooltip/index.js';
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
		target,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> & {
		numerator: number;
		denominator: number;
		ratio?: Snippet<[{ numerator: number; denominator: number }]>;
		detail?: Snippet<[{ numerator: number; denominator: number }]>;
		target: ProgressTargetType;
	} = $props();

	const progressRatio = $derived(denominator > 0 ? numerator / denominator : 0);
</script>

{#if denominator > 0}
	{#if ratio}
		<Tooltip.Provider>
			<Tooltip.Root>
				<Tooltip.Trigger>
					<div bind:this={ref} data-slot="progress-root" {...restProps}>
						<Progress
							value={progressRatio}
							max={1}
							class={formatProgressColor(numerator, denominator, target)}
						/>
						<div
							class={cn(
								'flex items-center justify-end font-light text-muted-foreground sm:min-w-[100px] md:min-w-[200px]',
								className
							)}
						>
							{@render ratio?.({ numerator, denominator })}
						</div>
					</div>
				</Tooltip.Trigger>
				<Tooltip.Content>
					{@render detail?.({ numerator, denominator })}
				</Tooltip.Content>
			</Tooltip.Root>
		</Tooltip.Provider>
	{:else}
		<Progress
			value={progressRatio}
			max={1}
			class={formatProgressColor(numerator, denominator, target)}
		/>
	{/if}
{:else}
	<Icon icon="ph:infinity" />
{/if}
