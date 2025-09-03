<script lang="ts">
	import { fade } from 'svelte/transition';
	import { SmoothStepEdge, type EdgeProps } from '@xyflow/svelte';
	import { cn } from '$lib/utils';

	let { selected, ...restProps }: EdgeProps = $props();
</script>

<g out:fade|global class={cn(selected ? 'animate-pulse' : 'animate-none')}>
	<SmoothStepEdge {...restProps} style="stroke-width: 0.5px; stroke-dasharray: none;" />
</g>
