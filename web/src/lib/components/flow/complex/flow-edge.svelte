<script lang="ts">
	import { BezierEdge, type EdgeProps } from '@xyflow/svelte';
	import { fade } from 'svelte/transition';

	import { cn } from '$lib/utils';

	let { selected, ...restProps }: EdgeProps = $props();
</script>

<g out:fade|global>
	<BezierEdge
		{...restProps}
		style={cn('stroke-dasharray: none;', selected ? 'stroke-width: 1px' : 'stroke-width: 0.5px')}
	/>
</g>
