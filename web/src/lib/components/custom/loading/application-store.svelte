<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
</script>

<script lang="ts">
	const lines = 2;
</script>

<div class="flex flex-col gap-4">
	<div class="flex items-center justify-center">
		<span class="*:h-[50px] *:w-[300px]">
			<Skeleton class="bg-muted" />
		</span>
	</div>
	<div class="flex items-center justify-center">
		<span class="*:h-[30px] *:w-[700px]">
			<Skeleton class="bg-muted" />
		</span>
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
	{#each Array(lines)}
		<div class="flex items-center gap-4 *:h-[350px] *:w-full">
			{#each Array(3)}
				<Skeleton class="bg-muted" />
			{/each}
		</div>
	{/each}
	<div class="flex items-center justify-center gap-2 *:h-[50px] *:w-[50px]">
		{#each Array(7)}
			<Skeleton class="bg-muted" />
		{/each}
	</div>
</div>
