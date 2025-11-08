<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<p class="my-2 max-w-7xl text-center text-sm text-muted-foreground">
	{@render children?.()}
</p>
