<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<div class="flex w-full max-w-7xl">
	<span class="ml-auto flex gap-2">
		{@render children?.()}
	</span>
</div>
