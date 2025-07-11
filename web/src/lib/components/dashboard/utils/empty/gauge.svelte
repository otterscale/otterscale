<script lang="ts" module>
	import { cn } from '$lib/utils.js';
	import { Arc, Chart, Group, Svg } from 'layerchart';

	const content = 'No Data';
	const data = Math.random() * 100;
</script>

<script lang="ts">
	let { class: className }: { class?: string } = $props();
</script>

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
		{content}
	</div>
</div>
