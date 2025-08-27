<script lang="ts" module>
	import Button from '$lib/components/ui/button/button.svelte';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { type WithElementRef } from 'bits-ui';
	import { getContext, type Snippet } from 'svelte';
	import type { HTMLAnchorAttributes, HTMLButtonAttributes } from 'svelte/elements';
	import { IndexManager, StepManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		text,
		icon,
		...restProps
	}: WithElementRef<HTMLButtonAttributes> &
		WithElementRef<HTMLAnchorAttributes> & {
			text?: Snippet<[{ value: number }]>;
			icon?: string;
		} = $props();

	const stepManager: StepManager = getContext('StepManager');
	const indexManager: IndexManager = getContext('IndexManager');
	const value = indexManager.get();
</script>

{#if value > 1}
	<div class="flex h-12 w-full items-center px-2">
		<Progress value={stepManager.areStepsActive(value) ? 100 : 0} class="bg-muted" />
	</div>
{/if}
<div class="flex flex-col items-center justify-center gap-1 transition-all duration-300">
	<div class="h-12 w-fit">
		<Button
			class={cn(
				'm-0 h-full w-full rounded-full border transition-all duration-300',
				stepManager.areStepsActive(value) ? 'bg-primary' : 'bg-muted',
				stepManager.isStepActive(value) && 'ring-card ring'
			)}
			onclick={async (e) => {
				e.preventDefault();
				await stepManager.update(value);
			}}
			{...restProps}
		>
			<Icon
				icon={icon ?? 'ph:empty'}
				class={cn(
					'size-6',
					icon ?? 'invisible',
					stepManager.areStepsActive(value) ? 'text-card' : 'text-muted-foreground'
				)}
			/>
		</Button>
	</div>

	{#if text}
		<p class="text-base whitespace-nowrap">{@render text({ value })}</p>
	{/if}
</div>
