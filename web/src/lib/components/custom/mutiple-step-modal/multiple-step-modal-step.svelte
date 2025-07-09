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
	<div class="flex h-10 w-full flex-col items-center justify-center p-2">
		<Progress value={stepManager.areStepsActive(value) ? 100 : 0} class="bg-muted" />
	</div>
{/if}

<div class="flex flex-col items-center justify-center gap-1 transition-all duration-300">
	{#if text}
		<p class="whitespace-nowrap text-base">{@render text({ value })}</p>
	{/if}

	<div class="h-10 w-10">
		<Button
			class={cn(
				'text-card m-0 h-10 w-10 rounded-full transition-all duration-300',
				stepManager.areStepsActive(value) ? 'bg-primary' : 'bg-muted',
				stepManager.isStepActive(value) && 'ring-primary ring-1 ring-offset-1'
			)}
			onclick={async (e) => {
				e.preventDefault();
				await stepManager.update(value);
			}}
			{...restProps}
		>
			<span>
				<Icon
					icon={icon ?? ''}
					class={cn(
						'size-6',
						stepManager.areStepsActive(value) ? 'text-card' : 'text-muted-foreground'
					)}
				/>
			</span>
		</Button>
	</div>
</div>
