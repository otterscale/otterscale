<script lang="ts">
	import { Badge } from '$lib/components/ui/badge';

	export let items: string[] = [];
</script>

<div class="grid gap-1">
	{#each items as item}
		<Badge variant="outline">
			{item}
		</Badge>
	{/each}
</div>
