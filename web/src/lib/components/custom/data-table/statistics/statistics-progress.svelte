<script lang="ts">
	import { Progress as ProgressPrimitive } from 'bits-ui';

	import {
		getProgressColor,
		type ProgressTargetType
	} from '$lib/components/custom/progress/utils.svelte';
	import * as Progress from '$lib/components/ui/progress/index.ts';
	import { cn, type WithoutChildrenOrChild } from '$lib/utils.js';

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
		target: ProgressTargetType;
	} = $props();
</script>

<Progress.Root
	bind:ref
	value={numerator / denominator}
	max={1}
	class={cn(
		getProgressColor(numerator, denominator, target),
		'absolute top-0 left-0 h-2 rounded-none',
		className
	)}
	{...restProps}
/>
