<script lang="ts" module>
	import * as Tabs from '$lib/components/custom/tabs/naive/index';
	import Progress from '$lib/components/ui/progress/progress.svelte';
	import { cn } from '$lib/utils.js';
	import Icon from '@iconify/svelte';
	import { Tabs as TabsPrimitive } from 'bits-ui';
	import { getContext } from 'svelte';
	import { IndexManager, StepManager } from './utils.svelte';
</script>

<script lang="ts">
	let {
		ref = $bindable(null),
		class: className,
		text,
		icon,
		...restProps
	}: Omit<TabsPrimitive.TriggerProps, 'value'> & {
		text?: string;
		icon?: string;
	} = $props();

	const stepManager: StepManager = getContext('StepManager');
	const indexManager: IndexManager = getContext('IndexManager');
	const value = indexManager.get();
</script>

{#if value > 1}
	<div class="flex h-10 w-full flex-col items-center justify-center p-2">
		<Progress value={stepManager.isStepActive(value) ? 100 : 0} class="bg-muted" />
	</div>
{/if}

<div class="flex flex-col items-center justify-center gap-2 transition-all duration-300">
	<p class="whitespace-nowrap text-base">{text ?? `Step ${value}`}</p>
	<div class="h-10 w-10">
		<Tabs.Trigger
			bind:ref
			data-slot="multiple-step-modal-step"
			value={String(value)}
			class={cn(
				'data-[state=active]:bg-primary data-[state=active]:ring-primary m-0 h-10 w-10 rounded-full transition-all duration-300 data-[state=active]:border-none data-[state=active]:ring-2 data-[state=active]:ring-offset-1',
				stepManager.isStepActive(value) ? 'bg-primary' : 'bg-muted'
			)}
			onclick={() => {
				stepManager.update(value);
			}}
			{...restProps}
		>
			<span>
				<Icon
					icon={icon ?? ''}
					class={cn('size-6', stepManager.isStepActive(value) && 'text-card')}
				/>
			</span>
		</Tabs.Trigger>
	</div>
</div>
