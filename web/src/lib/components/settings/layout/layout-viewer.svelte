<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<div class="my-4 flex w-full max-w-7xl items-center justify-center">
	{@render children?.()}
</div>
