<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<h1 class="text-lg font-bold">
	{@render children?.()}
</h1>
