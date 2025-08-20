<script lang="ts">
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
</script>

<div class="flex flex-col gap-4">
	<div class="flex flex-col gap-2">
		<Skeleton class="bg-muted h-[50px] w-[500px]" />
		<Skeleton class="bg-muted h-[20px] w-[100px]" />
		<div class="flex items-center gap-2 *:h-[20px] *:w-[100px]">
			{#each Array(3) as _}
				<Skeleton class="bg-muted" />
			{/each}
		</div>
	</div>

	<div class="flex items-center gap-2 *:h-[300px] *:w-[300px]">
		{#each Array(3) as _}
			<Skeleton class="bg-muted" />
		{/each}
	</div>

	{#each Array(3) as _}
		<div class="flex flex-col gap-2">
			<Skeleton class="bg-muted h-[50px] w-full " />
			{#each Array(3) as _}
				<div class="flex items-center gap-2 *:h-[30px] *:w-full">
					{#each Array(3) as _}
						<Skeleton class="bg-muted" />
					{/each}
				</div>
			{/each}
		</div>
	{/each}
</div>
