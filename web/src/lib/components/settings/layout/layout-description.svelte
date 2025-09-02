<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<p class="text-muted-foreground my-2 text-sm">
	{@render children?.()}
</p>
