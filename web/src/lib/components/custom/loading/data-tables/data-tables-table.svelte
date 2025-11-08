<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
</script>

<div class="flex flex-col gap-2">
	<div class="flex items-center gap-2 *:h-[300px] *:w-[300px]">
		{#each Array(3)}
			<Skeleton class="bg-muted" />
		{/each}
	</div>
	<div class="flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 *:h-[30px] *:w-[100px]">
			{#each Array(3)}
				<Skeleton class="bg-muted" />
			{/each}
		</div>
		<div class="flex items-center gap-2 *:h-[50px] *:w-[100px]">
			<Skeleton class="bg-muted" />
		</div>
	</div>
	<Skeleton class="h-[50px] w-full bg-muted " />
	{#each Array(12)}
		<div class="flex items-center gap-2 *:h-[30px] *:w-full">
			{#each Array(3)}
				<Skeleton class="bg-muted" />
			{/each}
		</div>
	{/each}
	<div class="flex items-center justify-between gap-2">
		<div class="flex items-center gap-2 *:h-[30px] *:w-[100px]">
			{#each Array(3)}
				<Skeleton class="bg-muted" />
			{/each}
		</div>
		<div class="flex items-center gap-2 *:h-[50px] *:w-[50px]">
			{#each Array(5)}
				<Skeleton class="bg-muted" />
			{/each}
		</div>
	</div>
</div>
