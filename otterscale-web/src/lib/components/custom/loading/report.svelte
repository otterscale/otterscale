<script lang="ts">
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
	import { cn } from '$lib/utils';

	let widthClasses: string[] = [];
	[2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 4, 4, 4, 5, 5, 6, 6, 12].forEach((denominator) => {
		for (let numerator = 1; numerator < denominator; numerator++) {
			widthClasses.push(`w-${numerator}/${denominator}`);
		}
	});
	function getWidthClass() {
		const widthClass = widthClasses[Math.floor(Math.random() * widthClasses.length)];
		return widthClass;
	}
</script>

<div class="flex flex-col gap-2">
	{#each Array(13) as _}
		<Skeleton class={cn('bg-muted h-[30px]', getWidthClass())} />
	{/each}
</div>
