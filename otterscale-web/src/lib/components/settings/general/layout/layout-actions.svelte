<script lang="ts" module>
	import type { Snippet } from 'svelte';
</script>

<script lang="ts">
	let { children }: { children?: Snippet } = $props();
</script>

<div class="flex">
	<span class="ml-auto space-x-2">
		{@render children?.()}
	</span>
</div>
