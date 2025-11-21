<script lang="ts">
	import { cn, type WithoutChildrenOrChild } from '$lib/utils.js';
	import { Progress as ProgressPrimitive } from 'bits-ui';
	import * as Progress from '$lib/components/ui/progress/index.ts';
	import { getProgressColor } from '$lib/components/custom/progress/utils.svelte';

	let {
		ref = $bindable(null),
		numerator,
		denominator,
		target,
		class: className,
		...restProps
	}: Omit<WithoutChildrenOrChild<ProgressPrimitive.RootProps>, 'value' | 'max'> & {
		numerator: number;
		denominator: number;
		target?: number;
	} = $props();
</script>

<Progress.Root
	bind:ref
	value={numerator / denominator}
	max={1}
	class={cn(
		getProgressColor(numerator, denominator, 'LTB'),
		'absolute top-0 left-0 h-2 rounded-none'
	)}
	{...restProps}
/>
