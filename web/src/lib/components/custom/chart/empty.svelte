<script lang="ts">
	import { cn } from '$lib/utils.js';
	import { curveCatmullRom } from 'd3-shape';
	import { AreaChart, Arc, Svg, Group, Chart, Text } from 'layerchart';

	let renderContext: 'svg' | 'canvas' = 'svg';
	let debug = false;

	const content = 'No Data';

	let {
		type = 'text',
		class: className
	}: { type?: 'text' | 'gauge' | 'area' | null | undefined; class?: string } = $props();
</script>

{#if type === 'text'}
	<p class={cn('flex h-full w-full items-center justify-center', className)}>{content}</p>
{:else if type === 'gauge'}
	{@const data = Math.random() * 100}
	<div class={cn('flex h-full w-full items-center justify-center', className)}>
		<div class="h-[200px] w-[200px]">
			<Chart>
				<Svg center>
					<Group>
						<Arc
							value={data}
							domain={[0, 100]}
							outerRadius={100}
							innerRadius={-13}
							cornerRadius={13}
							range={[-120, 120]}
							class={'fill-muted-foreground'}
							track={{ class: 'fill-muted' }}
						/>
					</Group>
				</Svg>
			</Chart>
		</div>
		<div class="absolute">
			<p>{content}</p>
		</div>
	</div>
{:else if type === 'area'}
	{@const data = Array.from({ length: 3 }, (_, index) => ({
		x: index + 1,
		y: Math.random()
	}))}
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
			<p>{content}</p>
		</div>
	</div>
{/if}
