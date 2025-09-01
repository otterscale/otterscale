<script lang="ts" module>
	import { Skeleton } from '$lib/components/ui/skeleton/index.js';
</script>

<script lang="ts">
	let { tabs = 2 }: { tabs?: number } = $props();
</script>

<div class="flex items-center gap-1 *:h-[30px] *:w-[100px]">
	{#each Array(tabs) as _}
		<Skeleton class="bg-muted first:bg-muted-foreground rounded-none first:rounded-l-lg last:rounded-r-lg" />
	{/each}
</div>
