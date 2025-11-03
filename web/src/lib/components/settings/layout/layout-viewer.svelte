<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<div class="flex items-center justify-center w-full my-4 max-w-7xl">
	{@render children?.()}
</div>
