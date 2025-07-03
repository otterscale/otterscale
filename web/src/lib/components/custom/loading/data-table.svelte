<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
</script>

<script lang="ts">
	let { lines = 10 }: { lines?: number } = $props();
</script>

<div class="flex flex-col gap-2">
	<div class="flex items-center gap-2 *:h-[200px] *:w-[200px]">
		<Skeleton class="bg-muted" />
		<Skeleton class="bg-muted" />
		<Skeleton class="bg-muted" />
	</div>
	<div class="flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 *:h-[30px] *:w-[100px]">
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
		</div>
		<div class="flex items-center gap-2 *:h-[30px] *:w-[100px]">
			<Skeleton class="bg-muted" />
		</div>
	</div>
	<Skeleton class="bg-muted h-[50px] w-full " />
	{#each Array(lines) as _}
		<div class="flex items-center gap-2 *:h-[20px] *:w-full">
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
		</div>
	{/each}
	<div class="flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 *:h-[20px] *:w-[100px]">
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
		</div>
		<div class="flex items-center gap-2 *:h-[30px] *:w-[30px]">
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
			<Skeleton class="bg-muted" />
		</div>
	</div>
</div>
