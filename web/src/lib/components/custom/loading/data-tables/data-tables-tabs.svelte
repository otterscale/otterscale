<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
</script>

<script lang="ts">
	let { tabs = 2 }: { tabs?: number } = $props();
</script>

<div class="flex items-center gap-1 *:h-[30px] *:w-[100px]">
	{#each Array(tabs)}
		<Skeleton
			class="rounded-none bg-muted first:rounded-l-lg first:bg-muted-foreground last:rounded-r-lg"
		/>
	{/each}
</div>
