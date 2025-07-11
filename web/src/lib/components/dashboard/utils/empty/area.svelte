<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import { curveCatmullRom } from 'd3-shape';
	import { AreaChart } from 'layerchart';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	const content = 'No Data';
	const data = Array.from({ length: 3 }, (_, index) => ({
		x: index + 1,
		y: Math.random()
	}));
</script>

<script lang="ts">
	let { class: className }: { class?: string } = $props();
</script>

<div class={cn('relative h-full w-full', className)}>
	<AreaChart
		{data}
		x="x"
		series={[{ key: 'y', color: 'hsl(var(--muted-foreground))' }]}
		tooltip={false}
		axis={false}
		props={{ area: { curve: curveCatmullRom } }}
		{renderContext}
		{debug}
	/>
	<div class="absolute inset-0 flex items-center justify-center">
		{content}
	</div>
</div>
