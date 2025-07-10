<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { cn } from '$lib/utils';

	let widthClasses: string[] = [];
	[2, 3, 4, 5, 6, 12].forEach((denominator) => {
		for (let numerator = 1; numerator < denominator; numerator++) {
			widthClasses.push(`w-${numerator}/${denominator}`);
		}
	});
	function getWidthClass() {
		const widthClass = widthClasses[Math.floor(Math.random() * widthClasses.length)];
		console.log(widthClass);
		return widthClass;
	}
</script>

<script lang="ts">
	let { lines = 10 }: { lines?: number } = $props();
</script>

<div class="flex flex-col gap-2">
	{#each Array(lines) as _}
		<Skeleton class={cn('bg-muted h-[20px]', getWidthClass())} />
	{/each}
</div>
