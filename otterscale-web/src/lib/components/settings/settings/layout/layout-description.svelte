<script lang="ts">
	import type { Snippet } from 'svelte';

	let { children }: { children?: Snippet } = $props();
</script>

<p class="text-muted-foreground text-sm">
	{@render children?.()}
</p>
